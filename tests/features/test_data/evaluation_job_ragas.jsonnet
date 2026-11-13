local test = import 'test.libsonnet';
local collectionId = test.value('collection_id');
// The embedding connection comes from embedding_url/embedding_api-key in the
// Secret referenced by MODEL_AUTH_SECRET_REF; the adapter uses sidecar refs.
local requiredEmbeddingEnv(name) =
  local value = std.stripChars(test.env(name), ' \t\r\n');
  assert value != '' : 'Set ' + name + ' to the embedding model ID before running Ragas tests';
  value;
local parameters = {
  num_examples: 1,
  embedding_model: requiredEmbeddingEnv('RAGAS_EMBEDDING_MODEL'),
};
local ragasBenchmark(id) =
  test.benchmark(id, 'ragas', parameters) +
  // Disconnected runs use the benchmark-scoped MinIO bundle assigned by
  // test.benchmark(). Connected runs retain the external RAGAS source.
  if test.isDisconnected() then {} else {
    test_data_ref: {
      s3: {
        bucket: 'ai-eng-canada',
        key: 'datasets/ragas/dataset.jsonl',
        secret_ref: 's3-ragas-credentials',
      },
    },
  };

test.mergeOptional(
  test.mergeOptional(
    {
      model: test.model(),
      name: 'test-evaluation-job-for-ragas-benchmark',
    } + if collectionId == '' then {
      benchmarks: [
        ragasBenchmark('ragas_rag_default'),
        ragasBenchmark('ragas_rag_full'),
      ],
      tags: ['benchmark-providers', 'ragas'],
    } else {},
    if collectionId != '' then test.collection() else null,
  ),
  test.experiment('my-test-experiment'),
)
