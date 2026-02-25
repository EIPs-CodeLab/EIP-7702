#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-dev}"
DIST_DIR="${2:-dist}"

EXAMPLES=(
  "basic-authorization"
  "transaction-batching"
  "send-userop"
)

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

for target in "${TARGETS[@]}"; do
  IFS='/' read -r GOOS GOARCH <<<"${target}"
  stage_dir="${DIST_DIR}/stage-${GOOS}-${GOARCH}"
  mkdir -p "${stage_dir}"

  for example in "${EXAMPLES[@]}"; do
    output_name="${example}"
    if [[ "${GOOS}" == "windows" ]]; then
      output_name+=".exe"
    fi

    CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
      go build -trimpath -ldflags "-s -w" \
      -o "${stage_dir}/${output_name}" "./examples/${example}"
  done

  archive_path="${DIST_DIR}/eip7702-go_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
  tar -C "${stage_dir}" -czf "${archive_path}" .
  rm -rf "${stage_dir}"
  echo "Built ${archive_path}"
done

if command -v sha256sum >/dev/null 2>&1; then
  (
    cd "${DIST_DIR}"
    sha256sum ./*.tar.gz > checksums.txt
  )
else
  (
    cd "${DIST_DIR}"
    shasum -a 256 ./*.tar.gz > checksums.txt
  )
fi

echo "Release artifacts are available in ${DIST_DIR}/"
