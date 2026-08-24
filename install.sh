#!/usr/bin/env sh
#
# Praxicraft Assess CLI installer
# https://github.com/praxicraft-platform/praxicraft-assess-cli
#
# Usage:
#   curl -fsSL https://praxicraft.com/install.sh | sh
#
# Environment overrides:
#   PRAXICRAFT_VERSION       Install a specific release tag (e.g. "v2.0.0"). Defaults to latest.
#   PRAXICRAFT_INSTALL_DIR   Directory for the `praxicraft-assess` binary. Auto-detected by default.
#   PRAXICRAFT_NO_VERIFY     Set to "1" to skip SHA-256 checksum verification (not recommended).
#
set -eu

REPO="praxicraft-platform/praxicraft-assess-cli"
DOCS_URL="https://docs.praxicraft.com/sdks/cli"
DEFAULT_VERSION="${PRAXICRAFT_VERSION:-}"
NO_VERIFY="${PRAXICRAFT_NO_VERIFY:-0}"
CMD_NAME="praxicraft-assess"

# ---- pretty output -----------------------------------------------------------

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  BOLD=$(printf '\033[1m')
  DIM=$(printf '\033[2m')
  RED=$(printf '\033[31m')
  GREEN=$(printf '\033[32m')
  YELLOW=$(printf '\033[33m')
  BLUE=$(printf '\033[34m')
  RESET=$(printf '\033[0m')
else
  BOLD=""; DIM=""; RED=""; GREEN=""; YELLOW=""; BLUE=""; RESET=""
fi

info()  { printf "%s==>%s %s\n" "${BLUE}${BOLD}" "${RESET}" "$*"; }
warn()  { printf "%swarning:%s %s\n" "${YELLOW}${BOLD}" "${RESET}" "$*" >&2; }
err()   { printf "%serror:%s %s\n" "${RED}${BOLD}" "${RESET}" "$*" >&2; }
ok()    { printf "%s✓%s %s\n" "${GREEN}${BOLD}" "${RESET}" "$*"; }

abort() {
  err "$1"
  printf "\nFor manual installation see: %s\n" "$DOCS_URL" >&2
  exit 1
}

# ---- prerequisites -----------------------------------------------------------

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || abort "missing required command: $1"
}

need_cmd uname
need_cmd mkdir
need_cmd mv
need_cmd chmod
need_cmd rm

if command -v curl >/dev/null 2>&1; then
  DOWNLOADER="curl"
elif command -v wget >/dev/null 2>&1; then
  DOWNLOADER="wget"
else
  abort "either 'curl' or 'wget' is required"
fi

fetch() {
  url="$1"
  out="$2"
  if [ "$DOWNLOADER" = "curl" ]; then
    curl -fsSL --proto '=https' --tlsv1.2 -o "$out" "$url"
  else
    wget -q -O "$out" "$url"
  fi
}

fetch_stdout() {
  url="$1"
  if [ "$DOWNLOADER" = "curl" ]; then
    curl -fsSL --proto '=https' --tlsv1.2 "$url"
  else
    wget -q -O - "$url"
  fi
}

# ---- platform detection ------------------------------------------------------

detect_os() {
  uname_s=$(uname -s 2>/dev/null || echo "")
  case "$uname_s" in
    Linux)   echo "linux" ;;
    Darwin)  echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*)
      abort "Windows is not supported by this installer. Download a binary from https://github.com/${REPO}/releases"
      ;;
    *) abort "unsupported operating system: $uname_s" ;;
  esac
}

detect_arch() {
  uname_m=$(uname -m 2>/dev/null || echo "")
  case "$uname_m" in
    x86_64|amd64) echo "x64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) abort "unsupported architecture: $uname_m" ;;
  esac
}

OS=$(detect_os)
ARCH=$(detect_arch)

# ---- resolve version ---------------------------------------------------------

resolve_version() {
  if [ -n "$DEFAULT_VERSION" ]; then
    echo "$DEFAULT_VERSION"
    return
  fi
  api="https://api.github.com/repos/${REPO}/releases/latest"
  tag=$(fetch_stdout "$api" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
  [ -n "$tag" ] || abort "could not resolve latest release tag from $api"
  echo "$tag"
}

VERSION=$(resolve_version)
BINARY="praxicraft-assess-${OS}-${ARCH}"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
BINARY_URL="${BASE_URL}/${BINARY}"
SUMS_URL="${BASE_URL}/SHA256SUMS.txt"

info "Installing Praxicraft Assess CLI ${BOLD}${VERSION}${RESET}${DIM} (${OS}/${ARCH})${RESET}"

# ---- pick install dir --------------------------------------------------------

pick_install_dir() {
  if [ -n "${PRAXICRAFT_INSTALL_DIR:-}" ]; then
    echo "$PRAXICRAFT_INSTALL_DIR"
    return
  fi
  for candidate in "/usr/local/bin" "$HOME/.local/bin" "$HOME/bin"; do
    if [ -d "$candidate" ] && [ -w "$candidate" ]; then
      echo "$candidate"
      return
    fi
  done
  echo "$HOME/.local/bin"
}

INSTALL_DIR=$(pick_install_dir)
mkdir -p "$INSTALL_DIR" || abort "cannot create $INSTALL_DIR"

# ---- download ----------------------------------------------------------------

TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t praxicraft-assess-cli)
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

info "Downloading ${DIM}${BINARY_URL}${RESET}"
fetch "$BINARY_URL" "$TMP_DIR/$BINARY" \
  || abort "failed to download $BINARY_URL"

# ---- checksum verification ---------------------------------------------------

verify_checksum() {
  sums_file="$TMP_DIR/SHA256SUMS.txt"

  if ! fetch "$SUMS_URL" "$sums_file"; then
    warn "could not download $SUMS_URL — skipping checksum verification"
    return 0
  fi

  expected_line=$(grep -E "[ *]${BINARY}$" "$sums_file" || true)
  if [ -z "$expected_line" ]; then
    warn "no checksum entry for $BINARY in SHA256SUMS.txt — skipping verification"
    return 0
  fi

  expected=$(echo "$expected_line" | awk '{print $1}')

  if command -v sha256sum >/dev/null 2>&1; then
    actual=$(sha256sum "$TMP_DIR/$BINARY" | awk '{print $1}')
  elif command -v shasum >/dev/null 2>&1; then
    actual=$(shasum -a 256 "$TMP_DIR/$BINARY" | awk '{print $1}')
  else
    warn "no sha256sum/shasum found — skipping checksum verification"
    return 0
  fi

  if [ "$expected" != "$actual" ]; then
    abort "checksum mismatch for $BINARY (expected $expected, got $actual)"
  fi
  ok "Checksum verified"
}

if [ "$NO_VERIFY" = "1" ]; then
  warn "PRAXICRAFT_NO_VERIFY=1 — skipping checksum verification"
else
  verify_checksum
fi

# ---- install -----------------------------------------------------------------

chmod +x "$TMP_DIR/$BINARY"
TARGET="$INSTALL_DIR/$CMD_NAME"

if [ -w "$INSTALL_DIR" ]; then
  mv -f "$TMP_DIR/$BINARY" "$TARGET"
else
  warn "$INSTALL_DIR is not writable; retrying with sudo"
  need_cmd sudo
  sudo mv -f "$TMP_DIR/$BINARY" "$TARGET"
fi

ok "Installed ${BOLD}${CMD_NAME}${RESET} to ${BOLD}$TARGET${RESET}"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    warn "$INSTALL_DIR is not on your PATH"
    printf "    Add it with one of:\n"
    printf "      echo 'export PATH=\"%s:\$PATH\"' >> ~/.bashrc\n" "$INSTALL_DIR"
    printf "      echo 'export PATH=\"%s:\$PATH\"' >> ~/.zshrc\n" "$INSTALL_DIR"
    ;;
esac

printf "\n"
ok "Done. Get started:"
printf "    %s%s configure%s   # save your API key\n" "$BOLD" "$CMD_NAME" "$RESET"
printf "    %s%s%s             # launch the interactive TUI\n" "$BOLD" "$CMD_NAME" "$RESET"
printf "    Docs: %s\n" "$DOCS_URL"
