local test = import 'test.libsonnet';

test.oobCollectionRefJobWithBenchmarks(
  'test-evaluation-job-oob-collection',
  'knowledge-reasoning-v1',
  [
    test.benchmark('gpqa:diamond', 'lighteval', { num_examples: 1, num_samples: 1 }),
    // HLE is judge-backed. Select the configured independent judge rather than
    // reusing the target model (the target for this suite is Qwen3-VL).
    // JUDGE_MODEL_NAME is supplied by the benchmark pipeline; do not hardcode
    // a model identity in the test fixture.
    test.benchmark('inspect/hle', 'inspect', {
      num_examples: 1,
      judge_model: test.env('JUDGE_MODEL_NAME'),
    }),
    test.benchmark('math_500', 'lighteval', { num_examples: 1, num_samples: 1 }),
  ],
)
