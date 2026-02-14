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

API_URL="https://api.github.com/repos/${REPO}/releases/latest"
RELEASE_JSON="$(curl -fsSL "$API_URL")"
TAG="$(printf "%s" "$RELEASE_JSON" | awk -F '\"' '/\"tag_name\":/ {print $4; exit}')"
[ -n "$TAG" ] || fail "Unable to determine latest release tag"

asset_exists() {
  printf "%s" "$RELEASE_JSON" | grep -q "\"name\": \"$1\""
}

ASSET_PRIMARY="${APP}-${OS}-${ARCH}.tar.gz"
ASSET_FALLBACK=""

if [ "$OS" = "darwin" ] && [ "$ARCH" = "arm64" ]; then
  ASSET_FALLBACK="${APP}-darwin-amd64.tar.gz"
fi

ASSET="$ASSET_PRIMARY"
if ! asset_exists "$ASSET_PRIMARY"; then
  if [ -n "$ASSET_FALLBACK" ] && asset_exists "$ASSET_FALLBACK"; then
    ASSET="$ASSET_FALLBACK"
    say "Asset ${ASSET_PRIMARY} not found in ${TAG}; using ${ASSET_FALLBACK}."
  else
    AVAILABLE_ASSETS="$(printf "%s" "$RELEASE_JSON" | sed -n 's/.*"name": "\([^"]*\)".*/\1/p' | grep "^${APP}-" | tr '\n' ' ')"
    fail "No installer asset for ${OS}/${ARCH} in ${TAG}. Available: ${AVAILABLE_ASSETS}"
  fi
fi

BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"
TMP_DIR="$(mktemp -d)"
ARCHIVE="${TMP_DIR}/${ASSET}"
CHECKSUMS="${TMP_DIR}/checksums.txt"
BIN_NAME="${ASSET%.tar.gz}"

say "Downloading ${ASSET} for ${OS}/${ARCH} (${TAG})..."
curl -fsSL -o "$ARCHIVE" "${BASE_URL}/${ASSET}"

EXPECTED=""
if curl -fsSL -o "$CHECKSUMS" "${BASE_URL}/checksums.txt" 2>/dev/null; then
  EXPECTED="$(grep " ${ASSET}\$" "$CHECKSUMS" | awk '{print $1}')"
else
  say "checksums.txt not found in release; using GitHub API digest."
fi

if [ -z "$EXPECTED" ]; then
  EXPECTED="$(printf "%s" "$RELEASE_JSON" | awk -v asset="$ASSET" '
    $0 ~ "\"name\": \"" asset "\"" { in_asset=1; next }
    in_asset && /"digest":/ {
      gsub(/[", ]/, "", $2)
      sub(/^sha256:/, "", $2)
      print $2
      exit
    }
    in_asset && /"name":/ { in_asset=0 }
  ')"
fi

[ -n "$EXPECTED" ] || fail "Checksum not found for ${ASSET}"

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL="$(sha256sum "$ARCHIVE" | awk '{print $1}')"
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL="$(shasum -a 256 "$ARCHIVE" | awk '{print $1}')"
else
  fail "Missing sha256 checker (sha256sum or shasum)"
fi
[ "$EXPECTED" = "$ACTUAL" ] || fail "Checksum mismatch for ${ASSET}"

tar -xzf "$ARCHIVE" -C "$TMP_DIR"

BIN_PATH="${TMP_DIR}/${BIN_NAME}"
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
