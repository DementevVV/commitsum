#!/bin/sh
set -eu

REPO="DementevVV/commitsum"
APP="commitsum"

say() {
  printf "%s\n" "$*"
}

fail() {
  say "Error: $*"
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Missing required command: $1"
}

need_cmd curl
need_cmd tar
need_cmd grep
need_cmd sed
need_cmd awk

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  mingw*|msys*|cygwin*) fail "Windows detected. Use install.ps1 instead." ;;
  *) fail "Unsupported OS: $OS" ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) fail "Unsupported architecture: $ARCH" ;;
esac

ASSET="${APP}-${OS}-${ARCH}.tar.gz"

API_URL="https://api.github.com/repos/${REPO}/releases/latest"
TAG="$(curl -fsSL "$API_URL" | awk -F '\"' '/\"tag_name\":/ {print $4; exit}')"
[ -n "$TAG" ] || fail "Unable to determine latest release tag"

BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"
TMP_DIR="$(mktemp -d)"
ARCHIVE="${TMP_DIR}/${ASSET}"
CHECKSUMS="${TMP_DIR}/checksums.txt"

say "Downloading ${ASSET} for ${OS}/${ARCH} (${TAG})..."
curl -fsSL -o "$ARCHIVE" "${BASE_URL}/${ASSET}"
curl -fsSL -o "$CHECKSUMS" "${BASE_URL}/checksums.txt"

if command -v sha256sum >/dev/null 2>&1; then
  CHECK_CMD="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
  CHECK_CMD="shasum -a 256"
else
  fail "Missing sha256 checker (sha256sum or shasum)"
fi

EXPECTED="$(grep " ${ASSET}\$" "$CHECKSUMS" | awk '{print $1}')"
[ -n "$EXPECTED" ] || fail "Checksum not found for ${ASSET}"

ACTUAL="$(eval "$CHECK_CMD" "$ARCHIVE" | awk '{print $1}')"
[ "$EXPECTED" = "$ACTUAL" ] || fail "Checksum mismatch for ${ASSET}"

tar -xzf "$ARCHIVE" -C "$TMP_DIR"

BIN_PATH="${TMP_DIR}/${APP}-${OS}-${ARCH}"
[ -f "$BIN_PATH" ] || fail "Binary not found in archive"

INSTALL_DIR="${HOME}/.local/bin"
if [ ! -d "$INSTALL_DIR" ]; then
  mkdir -p "$INSTALL_DIR"
fi

INSTALL_PATH="${INSTALL_DIR}/${APP}"
mv "$BIN_PATH" "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"

say "Installed ${APP} to ${INSTALL_PATH}"

case ":$PATH:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    say "Add this to your shell profile:"
    say "  export PATH=\"${INSTALL_DIR}:\\$PATH\""
    ;;
esac
