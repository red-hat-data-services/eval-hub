#!/usr/bin/env python3
"""Simple KFP pipeline test for Lighteval component.

This script creates a basic KFP pipeline that uses the Lighteval component
to demonstrate the integration.
"""

from kfp import compiler, dsl


@dsl.container_component
def lighteval_component(
    model_url: str,
    model_name: str,
    benchmark: str,
    tasks_json: str,
    num_fewshot: int = 0,
    limit: int = None,
    batch_size: int = 1,
) -> dsl.ContainerSpec:
    """Lighteval evaluation component.

    Args:
        model_url: Model endpoint URL
        model_name: Model identifier
        benchmark: Benchmark name
        tasks_json: JSON string of tasks to evaluate
        num_fewshot: Number of few-shot examples
        limit: Sample limit (optional)
        batch_size: Batch size

    Returns:
        Container specification
    """
    args_list = [
        "--model_url",
        model_url,
        "--model_name",
        model_name,
        "--benchmark",
        benchmark,
        "--tasks",
        tasks_json,
        "--num_fewshot",
        str(num_fewshot),
        "--batch_size",
        str(batch_size),
        "--output_metrics",
        "/tmp/outputs/metrics.json",
        "--output_results",
        "/tmp/outputs/results.json",
    ]

    if limit is not None:
        args_list.extend(["--limit", str(limit)])

    return dsl.ContainerSpec(
        image="quay.io/evalhub/lighteval-kfp:latest",
        command=["python", "/app/lighteval_component.py"],
        args=args_list,
    )


@dsl.pipeline(
    name="Lighteval Test Pipeline",
    description="Simple pipeline to test Lighteval component integration",
)
def lighteval_test_pipeline(
    model_url: str = "https://api.openai.com/v1",
    model_name: str = "gpt-3.5-turbo",
    benchmark: str = "mmlu",
    num_fewshot: int = 0,
    limit: int = 10,
):
    """Test pipeline for Lighteval component.

    Args:
        model_url: Model endpoint URL
        model_name: Model name/identifier
        benchmark: Benchmark to run
        num_fewshot: Number of few-shot examples
        limit: Limit number of samples for testing
    """
    # Define tasks to evaluate
    tasks_json = '["mmlu:abstract_algebra", "mmlu:anatomy"]'

    # Run lighteval component
    lighteval_task = lighteval_component(
        model_url=model_url,
        model_name=model_name,
        benchmark=benchmark,
        tasks_json=tasks_json,
        num_fewshot=num_fewshot,
        limit=limit,
        batch_size=1,
    )

    # Set display name for the component
    lighteval_task.set_display_name("Lighteval Evaluation")


def main():
    """Compile the pipeline to YAML."""
    pipeline_file = "lighteval_test_pipeline.yaml"

    compiler.Compiler().compile(
        pipeline_func=lighteval_test_pipeline, package_path=pipeline_file
    )

    print(f"✅ Pipeline compiled successfully: {pipeline_file}")
    print("\nTo run this pipeline:")
    print("1. Deploy a KFP instance (e.g., using Kind or OpenShift)")
    print("2. Upload the pipeline YAML to the KFP UI")
    print("3. Create a run with your model endpoint and API key")
    print("\nExample parameters:")
    print("  model_url: https://api.openai.com/v1")
    print("  model_name: gpt-3.5-turbo")
    print("  benchmark: mmlu")
    print("  num_fewshot: 0")
    print("  limit: 10")


if __name__ == "__main__":
    main()
