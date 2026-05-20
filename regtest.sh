#! /bin/bash

declare -a TEST=( "./..." )
if [ -n "$1" ]; then TEST=( "$@" ); fi

go build || exit 1

export REGTEST_BOT="${PWD}/sp0rkle"
export REGTEST_IRCD="${PWD}/ergo"

ERGO_VERSION="2.18.0"
ERGO_TAR="ergo-${ERGO_VERSION}-linux-x86_64.tar.gz"
ERGO_URL="https://github.com/ergochat/ergo/releases/download/v${ERGO_VERSION}/${ERGO_TAR}"
ERGO_TAR_SHA256="cbd888d9f89224eced6af76dae4b729eaa41ea04afd2e85fe9be8169a790a1da"
ERGO_BIN_SHA256="c50a20a999c787fe8b03e20988bd03d0641d47578de650e61e9a84876b06fa0d"

shasum_check() {
  local sum="$1"
  local file="$2"
  echo "${sum}  ${file}" | sha256sum --status --check -
}

ergo_tar_check() {
  shasum_check "${ERGO_TAR_SHA256}" "${ERGO_TAR}"
}

ergo_bin_check() {
  shasum_check "${ERGO_BIN_SHA256}" "${REGTEST_IRCD}"
}


if [ ! -f "${ERGO_TAR}" ] || ! ergo_tar_check; then
  wget -O "${ERGO_TAR}" "${ERGO_URL}"
  if ! ergo_tar_check; then
    echo "${ERGO_TAR}: checksum mismatch"
    exit 1
  fi
fi

if [ ! -f "${REGTEST_IRCD}" ] || ! ergo_bin_check; then
  tar --strip-components 1 -xzvf "${ERGO_TAR}" "ergo-${ERGO_VERSION}-linux-x86_64/ergo"
  if ! ergo_bin_check; then
    echo "${REGTEST_IRCD}: checksum mismatch"
    exit 1
  fi
fi

# -count 1 forces tests to be rerun every time instead of caching results
go test -v -count 1 -tags integration "${TEST[@]}"
