"""Lighteval framework schema adapter."""

import json
from typing import Any

from ...executors.base import ExecutionContext
from ...models.evaluation import EvaluationResult, EvaluationStatus
from ...utils.datetime_utils import safe_duration_seconds, utcnow
from ..base import SchemaAdapter
from ..transformers.benchmark import BenchmarkConfigTransformer
from ..transformers.metrics import MetricExtractor
from ..transformers.model import ModelConfigTransformer


class LightevalAdapter(SchemaAdapter):
    """Adapter for Lighteval framework (KFP-based).

    Lighteval is a lightweight evaluation framework from Hugging Face for
    evaluating language models. This adapter integrates Lighteval with
    eval-hub via Kubeflow Pipelines.

    The adapter generates KFP component specifications, transforms eval-hub
    execution contexts to Lighteval-compatible arguments, and parses
    Lighteval results back to eval-hub format.
    """

    def __init__(self) -> None:
        """Initialize Lighteval adapter."""
        super().__init__(framework_name="lighteval", version="1.0")
        self.model_transformer = ModelConfigTransformer()
        self.benchmark_transformer = BenchmarkConfigTransformer()
        self.metric_extractor = MetricExtractor()

    def get_kfp_component_spec(self) -> dict[str, Any]:
        """Generate KFP component specification for Lighteval.

        Returns:
            KFP component spec with inputs, outputs, and container implementation
        """
        return {
            "name": "lighteval-evaluate",
            "description": "Evaluate model using Lighteval framework",
            "inputs": [
                {
                    "name": "model_url",
                    "type": "String",
                    "description": "Model endpoint URL",
                },
                {
                    "name": "model_name",
                    "type": "String",
                    "description": "Model identifier",
                },
                {
                    "name": "benchmark",
                    "type": "String",
                    "description": "Benchmark name",
                },
                {
                    "name": "tasks",
                    "type": "JsonArray",
                    "description": "List of tasks to evaluate",
                },
                {
                    "name": "num_fewshot",
                    "type": "Integer",
                    "default": 0,
                    "description": "Number of few-shot examples",
                },
                {
                    "name": "limit",
                    "type": "Integer",
                    "optional": True,
                    "description": "Limit number of samples",
                },
                {
                    "name": "batch_size",
                    "type": "Integer",
                    "default": 1,
                    "description": "Batch size for evaluation",
                },
            ],
            "outputs": [
                {
                    "name": "output_metrics",
                    "type": "Metrics",
                    "description": "Evaluation metrics",
                },
                {
                    "name": "output_results",
                    "type": "Dataset",
                    "description": "Detailed results",
                },
            ],
            "implementation": {
                "container": {
                    "image": self.get_container_image(),
                    "command": ["python", "/app/lighteval_component.py"],
                    "args": [
                        "--model_url",
                        {"inputValue": "model_url"},
                        "--model_name",
                        {"inputValue": "model_name"},
                        "--benchmark",
                        {"inputValue": "benchmark"},
                        "--tasks",
                        {"inputValue": "tasks"},
                        "--num_fewshot",
                        {"inputValue": "num_fewshot"},
                        "--limit",
                        {"inputValue": "limit"},
                        "--batch_size",
                        {"inputValue": "batch_size"},
                        "--output_metrics",
                        {"outputPath": "output_metrics"},
                        "--output_results",
                        {"outputPath": "output_results"},
                    ],
                }
            },
        }

    def transform_to_kfp_args(
        self, context: ExecutionContext, backend_config: dict[str, Any]
    ) -> dict[str, Any]:
        """Transform eval-hub context to KFP component arguments.

        Args:
            context: Eval-hub execution context
            backend_config: Backend configuration

        Returns:
            Dictionary of arguments for Lighteval KFP component
        """
        # Use reusable transformers
        model_config = self.model_transformer.extract(context, backend_config)
        benchmark_config = self.benchmark_transformer.extract(
            context.benchmark_spec, backend_config
        )

        return {
            "model_url": model_config.url,
            "model_name": model_config.name,
            "benchmark": context.benchmark_spec.name,
            "tasks": context.benchmark_spec.tasks or [],
            "num_fewshot": benchmark_config.get("num_fewshot", 0),
            "limit": benchmark_config.get("limit"),
            "batch_size": benchmark_config.get("batch_size", 1),
        }

    def parse_kfp_output(
        self, artifacts: dict[str, str], context: ExecutionContext
    ) -> EvaluationResult:
        """Parse KFP component outputs to eval-hub result.

        Args:
            artifacts: Dictionary mapping artifact names to file paths
            context: Original execution context

        Returns:
            EvaluationResult with parsed metrics and metadata
        """
        # Load metrics from KFP artifact
        metrics = {}
        artifact_paths = {}

        if "output_metrics" in artifacts:
            metrics_path = artifacts["output_metrics"]
            artifact_paths["metrics"] = metrics_path

            try:
                with open(metrics_path) as f:
                    raw_metrics = json.load(f)
                    metrics = self.metric_extractor.extract(
                        raw_metrics,
                        framework="lighteval",
                        naming_strategy="hierarchical",
                    )
            except (FileNotFoundError, json.JSONDecodeError) as e:
                # If metrics file doesn't exist or is invalid, return empty metrics
                print(f"Warning: Could not load metrics from {metrics_path}: {e}")

        if "output_results" in artifacts:
            artifact_paths["results"] = artifacts["output_results"]

        completed_at = utcnow()
        duration = (
            safe_duration_seconds(completed_at, context.started_at)
            if context.started_at
            else 0.0
        )

        return EvaluationResult(
            evaluation_id=context.evaluation_id,
            provider_id="lighteval",
            benchmark_id=context.benchmark_spec.name,
            benchmark_name=context.benchmark_spec.name,
            status=EvaluationStatus.COMPLETED,
            metrics=dict(metrics),  # Ensure type compatibility
            artifacts=artifact_paths,
            error_message=None,
            started_at=context.started_at,
            completed_at=completed_at,
            duration_seconds=duration,
            mlflow_run_id=None,
        )

    def validate_config(self, config: dict[str, Any]) -> bool:
        """Validate configuration for Lighteval.

        Args:
            config: Backend configuration dictionary

        Returns:
            True if configuration is valid

        Raises:
            ValueError: If configuration is invalid
        """
        # Validate that framework is correctly specified
        if "framework" in config and config["framework"] != "lighteval":
            raise ValueError(
                f"Invalid framework '{config['framework']}' for LightevalAdapter. "
                "Expected 'lighteval'"
            )

        # Lighteval doesn't require additional config fields beyond what
        # KFP executor validates (kfp_endpoint, etc.)
        return True
