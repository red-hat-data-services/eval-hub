"""Functional tests for Lighteval adapter."""

import json
import tempfile
from datetime import UTC, datetime
from pathlib import Path

import pytest

from eval_hub.adapters.frameworks.lighteval import LightevalAdapter
from eval_hub.executors.base import ExecutionContext
from eval_hub.models.evaluation import BenchmarkSpec, EvaluationStatus


class TestLightevalAdapter:
    """Functional test suite for LightevalAdapter."""

    @pytest.fixture
    def adapter(self) -> LightevalAdapter:
        """Create a LightevalAdapter instance."""
        return LightevalAdapter()

    @pytest.fixture
    def execution_context(self) -> ExecutionContext:
        """Create a sample execution context."""
        from uuid import uuid4

        from eval_hub.models.evaluation import BackendSpec, BackendType

        benchmark_spec = BenchmarkSpec(
            name="mmlu",
            tasks=["mmlu_abstract_algebra", "mmlu_anatomy"],
            num_fewshot=5,
            batch_size=1,
            limit=10,
        )

        return ExecutionContext(
            evaluation_id=uuid4(),
            model_url="https://api.openai.com/v1",
            model_name="test-model",
            backend_spec=BackendSpec(
                name="kfp",
                type=BackendType.KFP,
                config={"framework": "lighteval"},
                benchmarks=[benchmark_spec],
            ),
            benchmark_spec=benchmark_spec,
            timeout_minutes=60,
            retry_attempts=3,
            started_at=datetime.now(UTC),
        )

    @pytest.fixture
    def backend_config(self) -> dict:
        """Create sample backend configuration."""
        return {
            "type": "kubeflow-pipelines",
            "framework": "lighteval",
            "model_url": "https://api.openai.com/v1",
        }

    def test_adapter_initialization(self, adapter: LightevalAdapter) -> None:
        """Test adapter initializes correctly."""
        assert adapter.framework_name == "lighteval"
        assert adapter.version == "1.0"

    def test_get_kfp_component_spec(self, adapter: LightevalAdapter) -> None:
        """Test KFP component specification is valid."""
        spec = adapter.get_kfp_component_spec()

        # Validate required structure
        assert spec["name"] == "lighteval-evaluate"
        assert "inputs" in spec
        assert "outputs" in spec
        assert "implementation" in spec

        # Validate required inputs exist
        input_names = {inp["name"] for inp in spec["inputs"]}
        required_inputs = {
            "model_url",
            "model_name",
            "benchmark",
            "tasks",
            "num_fewshot",
            "limit",
            "batch_size",
        }
        assert required_inputs.issubset(input_names)

        # Validate outputs exist
        output_names = {out["name"] for out in spec["outputs"]}
        assert {"output_metrics", "output_results"}.issubset(output_names)

    def test_transform_to_kfp_args(
        self,
        adapter: LightevalAdapter,
        execution_context: ExecutionContext,
        backend_config: dict,
    ) -> None:
        """Test context transformation produces valid KFP arguments."""
        args = adapter.transform_to_kfp_args(execution_context, backend_config)

        # Validate all required arguments are present
        assert args["model_url"] == backend_config["model_url"]
        assert args["model_name"] == execution_context.model_name
        assert args["benchmark"] == "mmlu"
        assert args["tasks"] == ["mmlu_abstract_algebra", "mmlu_anatomy"]
        assert args["num_fewshot"] == 5
        assert args["batch_size"] == 1
        assert args["limit"] == 10

    def test_parse_kfp_output_success(
        self, adapter: LightevalAdapter, execution_context: ExecutionContext
    ) -> None:
        """Test successful parsing of KFP outputs."""
        with tempfile.TemporaryDirectory() as tmpdir:
            metrics_path = Path(tmpdir) / "metrics.json"
            results_path = Path(tmpdir) / "results.json"

            # Write sample metrics and results
            metrics_data = {
                "mmlu_abstract_algebra.accuracy": 0.85,
                "mmlu_anatomy.accuracy": 0.78,
            }
            with open(metrics_path, "w") as f:
                json.dump(metrics_data, f)

            results_data = {
                "mmlu_abstract_algebra": {"accuracy": 0.85},
                "mmlu_anatomy": {"accuracy": 0.78},
            }
            with open(results_path, "w") as f:
                json.dump(results_data, f)

            artifacts = {
                "output_metrics": str(metrics_path),
                "output_results": str(results_path),
            }

            result = adapter.parse_kfp_output(artifacts, execution_context)

            # Validate result structure
            assert result.evaluation_id == execution_context.evaluation_id
            assert result.provider_id == "lighteval"
            assert result.benchmark_id == "mmlu"
            assert result.status == EvaluationStatus.COMPLETED
            assert len(result.metrics) > 0
            assert result.metrics["mmlu_abstract_algebra.accuracy"] == 0.85

    def test_parse_kfp_output_missing_files(
        self, adapter: LightevalAdapter, execution_context: ExecutionContext
    ) -> None:
        """Test parsing handles missing artifact files gracefully."""
        artifacts = {"output_metrics": "/nonexistent/metrics.json"}

        result = adapter.parse_kfp_output(artifacts, execution_context)

        # Should return result with empty metrics (no crash)
        assert result.evaluation_id == execution_context.evaluation_id
        assert result.status == EvaluationStatus.COMPLETED
        assert len(result.metrics) == 0

    def test_validate_config(self, adapter: LightevalAdapter) -> None:
        """Test config validation accepts valid configs."""
        assert adapter.validate_config({}) is True
        assert adapter.validate_config({"framework": "lighteval"}) is True

    def test_validate_config_rejects_wrong_framework(
        self, adapter: LightevalAdapter
    ) -> None:
        """Test config validation rejects incorrect framework."""
        with pytest.raises(ValueError, match="Invalid framework"):
            adapter.validate_config({"framework": "wrong-framework"})
