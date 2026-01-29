#!/bin/bash
# Build script for Aliasly
# Creates binaries for macOS and Linux (amd64 and arm64)

set -e

# Get version from git tag or use default
VERSION=${VERSION:-"0.1.4"}

# Output directory
OUT_DIR="dist"

# Clean and create output directory
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

echo "Building Aliasly v${VERSION}..."
echo ""

# Build for each platform
platforms=(
    "darwin/amd64"    # macOS Intel
    "darwin/arm64"    # macOS Apple Silicon
    "linux/amd64"     # Linux x86_64
    "linux/arm64"     # Linux ARM64
)

for platform in "${platforms[@]}"; do
    GOOS="${platform%/*}"
    GOARCH="${platform#*/}"
    output_name="al-${GOOS}-${GOARCH}"

    echo "Building ${output_name}..."

    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X aliasly/cmd.Version=${VERSION}" \
        -o "${OUT_DIR}/${output_name}" \
        .

    # Create compressed archive
    cd "$OUT_DIR"
    if [ "$GOOS" = "darwin" ]; then
        zip -q "${output_name}.zip" "${output_name}"
    else
        tar -czf "${output_name}.tar.gz" "${output_name}"
    fi
    cd ..
done

echo ""
echo "Build complete! Binaries are in ${OUT_DIR}/"
ls -lh "$OUT_DIR"
