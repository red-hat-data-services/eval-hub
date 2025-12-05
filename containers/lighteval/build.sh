#!/bin/bash
# Build script for Lighteval KFP component container

set -e

# Configuration
IMAGE_NAME="quay.io/evalhub/lighteval-kfp"
VERSION="${VERSION:-latest}"
PLATFORM="${PLATFORM:-linux/amd64}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Lighteval KFP Component Container${NC}"
echo "Image: ${IMAGE_NAME}:${VERSION}"
echo "Platform: ${PLATFORM}"
echo ""

# Change to script directory
cd "$(dirname "$0")"

# Build container
echo -e "${YELLOW}Building container image...${NC}"
docker build \
    --platform "${PLATFORM}" \
    -t "${IMAGE_NAME}:${VERSION}" \
    -t "${IMAGE_NAME}:latest" \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Container built successfully${NC}"
else
    echo -e "${RED}❌ Container build failed${NC}"
    exit 1
fi

# Optional: Push to registry
if [ "$PUSH" = "true" ]; then
    echo -e "${YELLOW}Pushing container to registry...${NC}"
    docker push "${IMAGE_NAME}:${VERSION}"
    docker push "${IMAGE_NAME}:latest"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Container pushed successfully${NC}"
    else
        echo -e "${RED}❌ Container push failed${NC}"
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}Build complete!${NC}"
echo "To run locally:"
echo "  docker run -it ${IMAGE_NAME}:${VERSION} --help"
echo ""
echo "To push to registry:"
echo "  PUSH=true ./build.sh"
