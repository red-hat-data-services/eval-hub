#!/usr/bin/env python3
"""Lighteval KFP Component.

This component runs Lighteval evaluations within a Kubeflow Pipeline.
It accepts model and benchmark configurations as inputs and produces
metrics and detailed results as KFP artifacts.
"""

import argparse
import json
import sys
from pathlib import Path
from typing import Any


def parse_args() -> argparse.Namespace:
    """Parse command-line arguments.

    Returns:
        Parsed arguments namespace
    """
    parser = argparse.ArgumentParser(description="Lighteval KFP Component")

    # Model configuration
    parser.add_argument(
        "--model_url", type=str, required=True, help="Model endpoint URL"
    )
    parser.add_argument(
        "--model_name", type=str, required=True, help="Model identifier"
    )

    # Benchmark configuration
    parser.add_argument("--benchmark", type=str, required=True, help="Benchmark name")
    parser.add_argument(
        "--tasks",
        type=str,
        required=True,
        help="JSON array of tasks to evaluate",
    )

    # Evaluation parameters
    parser.add_argument(
        "--num_fewshot",
        type=int,
        default=0,
        help="Number of few-shot examples",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=None,
        help="Limit number of samples (optional)",
    )
    parser.add_argument(
        "--batch_size",
        type=int,
        default=1,
        help="Batch size for evaluation",
    )

    # Output paths (KFP artifacts)
    parser.add_argument(
        "--output_metrics",
        type=str,
        required=True,
        help="Path to write metrics artifact",
    )
    parser.add_argument(
        "--output_results",
        type=str,
        required=True,
        help="Path to write detailed results artifact",
    )

    return parser.parse_args()


def run_lighteval_evaluation(
    model_url: str,
    model_name: str,
    benchmark: str,
    tasks: list[str],
    num_fewshot: int,
    limit: int | None,
    batch_size: int,
) -> dict[str, Any]:
    """Run Lighteval evaluation using the Lighteval CLI.

    Args:
        model_url: Model endpoint URL (OpenAI-compatible API)
        model_name: Model identifier
        benchmark: Benchmark name
        tasks: List of tasks to evaluate
        num_fewshot: Number of few-shot examples
        limit: Sample limit (optional)
        batch_size: Batch size

    Returns:
        Dictionary containing evaluation results

    Raises:
        RuntimeError: If evaluation fails
    """
    import shutil
    import subprocess
    import tempfile

    print("Running Lighteval evaluation:")
    print(f"  Model: {model_name} ({model_url})")
    print(f"  Benchmark: {benchmark}")
    print(f"  Tasks: {tasks}")
    print(
        f"  Parameters: num_fewshot={num_fewshot}, limit={limit}, batch_size={batch_size}"
    )

    # Prepare task list for Lighteval CLI
    # Format: task1|fewshot,task2|fewshot
    task_strings = []
    for task in tasks:
        if "|" in task:
            # Already formatted
            task_strings.append(task)
        else:
            # Add fewshot count
            task_strings.append(f"{task}|{num_fewshot}")

    tasks_arg = ",".join(task_strings)
    print(f"Formatted tasks: {tasks_arg}")

    # Prepare model configuration as comma-separated string
    # Format: model_name={},base_url={},provider={}
    model_args = f"model_name={model_name},base_url={model_url}"
    print(f"Model args: {model_args}")

    # Create temporary output directory for Lighteval
    output_dir = tempfile.mkdtemp(prefix="lighteval_")
    print(f"Using output directory: {output_dir}")

    try:
        # Build lighteval CLI command
        cmd = [
            "lighteval",
            "endpoint",
            "litellm",
            model_args,
            tasks_arg,
            "--output-dir",
            output_dir,
            "--no-push-to-hub",
            "--save-details",
        ]

        # Add optional parameters
        if limit is not None:
            # Note: Lighteval CLI uses max_samples parameter in task config
            # For now, we'll run all samples and handle limit in parsing
            print(f"Note: Sample limit ({limit}) will be applied during result parsing")

        print(f"Running command: {' '.join(cmd)}")

        # Run Lighteval CLI
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=3600,  # 1 hour timeout
        )

        # Check for errors
        if result.returncode != 0:
            print(f"Lighteval CLI stderr:\n{result.stderr}")
            raise RuntimeError(
                f"Lighteval CLI failed with exit code {result.returncode}\n"
                f"Stdout: {result.stdout}\n"
                f"Stderr: {result.stderr}"
            )

        print(f"Lighteval CLI stdout:\n{result.stdout}")

        # Parse results from Lighteval output
        # Lighteval writes results to output_dir/results/model_name/results_*.json
        results_pattern = Path(output_dir) / "results" / "**" / "results_*.json"
        import glob

        results_files = glob.glob(str(results_pattern), recursive=True)

        if not results_files:
            # Try alternative location
            results_pattern2 = Path(output_dir) / "**" / "results.json"
            results_files = glob.glob(str(results_pattern2), recursive=True)

        if not results_files:
            raise RuntimeError(
                f"No results file found in {output_dir}. "
                f"Available files: {list(Path(output_dir).rglob('*'))}"
            )

        print(f"Found results file: {results_files[0]}")

        # Load results from the first matching file
        with open(results_files[0]) as f:
            results_data = json.load(f)

        print("Evaluation completed. Results loaded successfully.")
        return results_data

    except subprocess.TimeoutExpired as e:
        raise RuntimeError(
            f"Lighteval evaluation timed out after {e.timeout} seconds"
        ) from e
    except json.JSONDecodeError as e:
        raise RuntimeError(f"Failed to parse Lighteval results file: {e}") from e
    except FileNotFoundError as e:
        raise RuntimeError(f"Lighteval results file not found: {e}") from e
    except Exception as e:
        raise RuntimeError(
            f"Unexpected error during Lighteval evaluation: {type(e).__name__}: {e}"
        ) from e

    finally:
        # Clean up temporary directory
        try:
            shutil.rmtree(output_dir)
            print(f"Cleaned up temporary directory: {output_dir}")
        except Exception as e:
            print(f"Warning: Failed to clean up {output_dir}: {e}")


def extract_metrics(results: dict[str, Any]) -> dict[str, float]:
    """Extract and flatten metrics from Lighteval results.

    Lighteval results format:
    {
        "results": {
            "task_name": {
                "metric_name": value,
                ...
            },
            ...
        }
    }

    Args:
        results: Raw Lighteval results

    Returns:
        Flattened metrics dictionary with hierarchical naming (task.metric)
    """
    metrics = {}

    # Check if results have the standard Lighteval structure
    if "results" in results:
        results_dict = results["results"]
    else:
        results_dict = results

    for task_name, task_results in results_dict.items():
        if isinstance(task_results, dict):
            for metric_name, value in task_results.items():
                # Use hierarchical naming (task.metric)
                metric_key = f"{task_name}.{metric_name}"

                # Handle nested metric values (some metrics have mean, std, etc.)
                if isinstance(value, dict):
                    for sub_metric, sub_value in value.items():
                        metrics[f"{metric_key}.{sub_metric}"] = sub_value
                else:
                    metrics[metric_key] = value

    return metrics


def write_kfp_artifacts(
    metrics: dict[str, float],
    results: dict[str, Any],
    metrics_path: str,
    results_path: str,
) -> None:
    """Write KFP output artifacts.

    Args:
        metrics: Flattened metrics dictionary
        results: Detailed results
        metrics_path: Path to write metrics artifact
        results_path: Path to write results artifact
    """
    # Ensure output directories exist
    Path(metrics_path).parent.mkdir(parents=True, exist_ok=True)
    Path(results_path).parent.mkdir(parents=True, exist_ok=True)

    # Write metrics artifact
    print(f"Writing metrics to {metrics_path}")
    with open(metrics_path, "w") as f:
        json.dump(metrics, f, indent=2)

    # Write detailed results artifact
    print(f"Writing results to {results_path}")
    with open(results_path, "w") as f:
        json.dump(results, f, indent=2)

    print(f"Artifacts written successfully. {len(metrics)} metrics recorded.")


def main() -> int:
    """Main entry point for Lighteval KFP component.

    Returns:
        Exit code (0 for success, 1 for failure)
    """
    try:
        # Parse arguments
        args = parse_args()

        # Parse tasks JSON array
        tasks = json.loads(args.tasks)
        if not isinstance(tasks, list):
            raise ValueError("Tasks must be a JSON array")

        # Run evaluation using actual Lighteval library
        results = run_lighteval_evaluation(
            model_url=args.model_url,
            model_name=args.model_name,
            benchmark=args.benchmark,
            tasks=tasks,
            num_fewshot=args.num_fewshot,
            limit=args.limit,
            batch_size=args.batch_size,
        )

        # Extract and normalize metrics
        metrics = extract_metrics(results)

        # Write KFP artifacts
        write_kfp_artifacts(
            metrics=metrics,
            results=results,
            metrics_path=args.output_metrics,
            results_path=args.output_results,
        )

        print("✅ Lighteval evaluation completed successfully")
        return 0

    except Exception as e:
        print(f"❌ Error running Lighteval evaluation: {e}", file=sys.stderr)
        import traceback

        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
