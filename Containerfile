# Multi-stage build for the evaluation hub
FROM registry.access.redhat.com/ubi9/python-312-minimal:latest as builder

USER 0

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Set work directory
WORKDIR /app

# Copy dependency files first for better caching
COPY pyproject.toml README.md requirements.txt ./

# Install dependencies (compatible with hermetic/cachi2 builds)
RUN . /cachi2/cachi2.env && \
    pip install -r requirements.txt && \
    pip install --no-deps -e .

# Copy source code after dependencies are installed
COPY src/ ./src/

# Production stage
FROM registry.access.redhat.com/ubi9/python-312-minimal:latest as production

USER 0

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Set work directory
WORKDIR /app

# Copy application code from builder stage (using numeric UID 1001 for UBI9 default user)
COPY --from=builder --chown=1001:0 /app/src ./src
COPY --from=builder --chown=1001:0 /app/pyproject.toml /app/README.md /app/requirements.txt ./

# Install the package in production stage (compatible with hermetic/cachi2 builds)
RUN . /cachi2/cachi2.env && \
    pip install -r requirements.txt && \
    pip install --no-deps -e .

# Create required directories and set permissions
RUN mkdir -p /app/logs /app/temp && \
    chmod -R g=u /app && \
    chmod 755 /app/src/eval_hub/data && \
    chmod 644 /app/src/eval_hub/data/providers.yaml

# Switch to non-root user (UID 1001 is the default user in UBI9 Python images)
USER 1001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/api/v1/health || exit 1

# Expose port
EXPOSE 8000

# Run the application
CMD ["/opt/app-root/bin/python3", "-m", "eval_hub.main"]
