# KFP Component Containers

This directory contains Kubeflow Pipelines (KFP) component implementations for various evaluation frameworks.

## Structure

```
containers/
├── base/
│   └── component_template.py    # Template for creating new components
└── lighteval/
    ├── Dockerfile                # Container image definition
    ├── lighteval_component.py    # Component implementation
    ├── run_component.sh          # Entrypoint script
    ...
    └── build.sh                  # Build automation script
```
