#!/usr/bin/env python3
"""Base template for KFP evaluation components.

This template provides a standard structure for creating new evaluation
framework components. Copy and customize for specific frameworks.
"""

import argparse
import json
import sys
from pathlib import Path
from typing import Any


def parse_args() -> argparse.Namespace:
    """Parse command-line arguments.

    Customize this function to add framework-specific arguments.

    Returns:
        Parsed arguments namespace
    """
    parser = argparse.ArgumentParser(description="Evaluation KFP Component Template")

    # Standard model configuration arguments
    parser.add_argument(
        "--model_url", type=str, required=True, help="Model endpoint URL"
    )
    parser.add_argument(
        "--model_name", type=str, required=True, help="Model identifier"
    )

    # Standard benchmark configuration arguments
    parser.add_argument("--benchmark", type=str, required=True, help="Benchmark name")
    parser.add_argument(
        "--tasks",
        type=str,
        required=True,
        help="JSON array of tasks to evaluate",
    )

    # Standard evaluation parameters
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

    # Add framework-specific arguments here
    # parser.add_argument("--framework_arg", ...)

    # Standard KFP output paths
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


def run_evaluation(
    model_url: str,
    model_name: str,
    benchmark: str,
    tasks: list[str],
    num_fewshot: int,
    limit: int | None,
    batch_size: int,
    **kwargs: Any,
) -> dict[str, Any]:
    """Run framework-specific evaluation.

    Customize this function to integrate with your evaluation framework.

    Args:
        model_url: Model endpoint URL
        model_name: Model identifier
        benchmark: Benchmark name
        tasks: List of tasks to evaluate
        num_fewshot: Number of few-shot examples
        limit: Sample limit (optional)
        batch_size: Batch size
        **kwargs: Additional framework-specific arguments

    Returns:
        Dictionary containing evaluation results

    Raises:
        ImportError: If framework library is not installed
        RuntimeError: If evaluation fails
    """
    # TODO: Import your evaluation framework
    # from your_framework import Evaluator

    print("Running evaluation:")
    print(f"  Model: {model_name} ({model_url})")
    print(f"  Benchmark: {benchmark}")
    print(f"  Tasks: {tasks}")
    print(
        f"  Parameters: num_fewshot={num_fewshot}, limit={limit}, batch_size={batch_size}"
    )

    # TODO: Configure and run your evaluator
    # evaluator = Evaluator(
    #     model=model_name,
    #     model_endpoint=model_url,
    #     tasks=tasks,
    #     num_fewshot=num_fewshot,
    #     limit=limit,
    #     batch_size=batch_size,
    # )
    # results = evaluator.evaluate()

    # Placeholder results - replace with actual evaluation
    results = {}
    for task in tasks:
        results[task] = {
            "accuracy": 0.75,
            "f1": 0.72,
        }

    print(f"Evaluation completed. {len(results)} tasks evaluated.")
    return results


def extract_metrics(results: dict[str, Any]) -> dict[str, float]:
    """Extract and flatten metrics from framework results.

    Customize this function to handle your framework's result format.

    Args:
        results: Raw evaluation results

    Returns:
        Flattened metrics dictionary with hierarchical naming (task.metric)
    """
    metrics = {}

    # Handle framework-specific result structure
    # Common pattern: iterate over tasks and metrics
    for task_name, task_results in results.items():
        if isinstance(task_results, dict):
            for metric_name, value in task_results.items():
                # Use hierarchical naming: task.metric
                metric_key = f"{task_name}.{metric_name}"

                # Handle nested values if needed
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

    This function is standard across all components.

    Args:
        metrics: Flattened metrics dictionary
        results: Detailed results
        metrics_path: Path to write metrics artifact
        results_path: Path to write results artifact
    """
    # Ensure output directories exist
    Path(metrics_path).parent.mkdir(parents=True, exist_ok=True)
    Path(results_path).parent.mkdir(parents=True, exist_ok=True)

    # Write metrics artifact (JSON format)
    print(f"Writing metrics to {metrics_path}")
    with open(metrics_path, "w") as f:
        json.dump(metrics, f, indent=2)

    # Write detailed results artifact (JSON format)
    print(f"Writing results to {results_path}")
    with open(results_path, "w") as f:
        json.dump(results, f, indent=2)

    print(f"Artifacts written successfully. {len(metrics)} metrics recorded.")


def main() -> int:
    """Main entry point for KFP component.

    This function is standard across all components.

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

        # Extract framework-specific args (customize as needed)
        framework_args = {}
        # framework_args["framework_arg"] = args.framework_arg

        # Run evaluation
        results = run_evaluation(
            model_url=args.model_url,
            model_name=args.model_name,
            benchmark=args.benchmark,
            tasks=tasks,
            num_fewshot=args.num_fewshot,
            limit=args.limit,
            batch_size=args.batch_size,
            **framework_args,
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

        print("✅ Evaluation completed successfully")
        return 0

    except Exception as e:
        print(f"❌ Error running evaluation: {e}", file=sys.stderr)
        import traceback

        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
