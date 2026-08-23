#!/usr/bin/env sh
# Install praxicraft-assess from the latest GitHub Release.
# Usage: curl -fsSL https://raw.githubusercontent.com/praxicraft-platform/praxicraft-assess-cli/main/install.sh | sh
set -eu

REPO="praxicraft-platform/praxicraft-assess-cli"
BINARY="praxicraft-assess"
INSTALL_DIR="${PRAXICRAFT_INSTALL_DIR:-}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  linux|darwin) ;;
  msys*|cygwin*|mingw*)
    echo "Windows: download a .zip from https://github.com/${REPO}/releases/latest" >&2
    exit 1
    ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2
  exit 1
fi

api="https://api.github.com/repos/${REPO}/releases/latest"
tag="$(curl -fsSL "$api" | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
if [ -z "$tag" ]; then
  echo "Could not resolve latest release for ${REPO}" >&2
  exit 1
fi

version="${tag#v}"
asset="${BINARY}_${version}_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${tag}/${asset}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT HUP

echo "Downloading ${BINARY} ${tag} (${os}/${arch})…"
curl -fsSL --progress-bar -o "${tmpdir}/${asset}" "$url"
tar -xzf "${tmpdir}/${asset}" -C "$tmpdir"

bin_path="$(find "$tmpdir" -type f -name "$BINARY" | head -n 1)"
if [ -z "$bin_path" ] || [ ! -f "$bin_path" ]; then
  echo "Archive did not contain ${BINARY}" >&2
  exit 1
fi
chmod +x "$bin_path"

if [ -z "$INSTALL_DIR" ]; then
  if [ -w /usr/local/bin ]; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="${HOME}/.local/bin"
  fi
fi
mkdir -p "$INSTALL_DIR"
dest="${INSTALL_DIR}/${BINARY}"
mv "$bin_path" "$dest"

echo "Installed ${dest}"
if ! command -v "$BINARY" >/dev/null 2>&1; then
  echo "Add to PATH: export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
"$dest" version 2>/dev/null || true
echo "Next: praxicraft-assess configure"
