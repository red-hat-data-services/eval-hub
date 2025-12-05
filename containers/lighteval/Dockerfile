FROM registry.redhat.io/ubi9/python-312:latest

# Set working directory
WORKDIR /app

# Upgrade pip
RUN pip install --no-cache-dir --upgrade pip

# Install Python dependencies
# Install lighteval and KFP SDK
RUN pip install --no-cache-dir \
    lighteval>=0.4.0 \
    kfp>=2.7.0 \
    openai>=1.0.0 \
    litellm>=1.0.0 \
    transformers>=4.35.0 \
    accelerate>=0.25.0 \
    torch>=2.1.0

# Copy component script
COPY lighteval_component.py /app/

# Set entrypoint
ENTRYPOINT ["python", "/app/lighteval_component.py"]
