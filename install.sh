#!/usr/bin/env bash

# This install script is intended to download and install the latest available
# release of the ocictl from ocibuilder.
#
# It attempts to identify the current platform and an error will be thrown if
# the platform is not supported.
#
# Environment variables:
# - INSTALL_DIRECTORY (optional): defaults to $GOPATH/bin
#
# You can install using this script:
# $ curl https://raw.githubusercontent.com/ocibuilder/ocibuilder/master/install.sh | sh

set -e

RELEASES_URL="https://github.com/ocibuilder/ocibuilder/releases"
LATEST_VERSION=$(curl --silent "$RELEASES_URL/latest" | sed 's#.*tag/\(.*\)\".*#\1#')
DOWNLOAD_BIN="ocictl/ocictl"

downloadTar() {
    url="$DOWNLOAD_URL"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        curl -0L $url | tar -xz
    elif test -x "$(command -v wget)"; then
        wget -c $url -O - | tar -xz
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi
    echo "Finished downloading tar"
}

findGoBinDirectory() {
    EFFECTIVE_GOPATH=$(go env GOPATH)
    if [ -z "$EFFECTIVE_GOPATH" ]; then
        echo "Installation could not determine your \$GOPATH."
        exit 1
    fi
    if [ -z "$GOBIN" ]; then
        GOBIN=$(echo "${EFFECTIVE_GOPATH%%:*}/bin" | sed s#//*#/#g)
    fi
    if [ ! -d "$GOBIN" ]; then
        echo "Installation requires your GOBIN directory $GOBIN to exist. Please create it."
        exit 1
    fi
    eval "$1='$GOBIN'"
}

initOS() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    OS_CYGWIN=0
    if [ -n "$OCI_OS" ]; then
        echo "Using OCI_OS"
        OS="$OCI_OS"
    fi
    case "$OS" in
        darwin) OS='darwin';;
        linux) OS='linux';;
        *) echo "OS ${OS} is not supported by this installation script"; exit 1;;
    esac
    echo "OS = $OS"
}

initOS

# determine install directory if required
if [ -z "$INSTALL_DIRECTORY" ]; then
    findGoBinDirectory INSTALL_DIRECTORY
fi
echo "Will install into $INSTALL_DIRECTORY"

if [ $OS = "darwin" ]; then
    DOWNLOAD_URL="$RELEASES_URL/download/$LATEST_VERSION/ocictl-darwin-amd64.tar.gz"
fi

if [ $OS = "linux" ]; then
    DOWNLOAD_URL="$RELEASES_URL/download/$LATEST_VERSION/ocictl-linux-amd64.tar.gz"
fi
echo "Downloading from $DOWNLOAD_URL"

downloadTar

echo "Setting executable permissions."
chmod +x "$DOWNLOAD_BIN"

echo "Moving executable to $INSTALL_DIRECTORY/$INSTALL_NAME"
mv "$DOWNLOAD_BIN" "$INSTALL_DIRECTORY"
