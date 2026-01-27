#!/bin/bash
OS=$(uname | tr '[:upper:]' '[:lower:]')
echo "Detected OS: $OS"

case "$OS" in
  linux)
    ./build_linux.sh
    ;;
  darwin)
    ./build_mac.sh
    ;;
  *)
    echo "Unsupported OS. Please run the corresponding build script."
    ;;
esac
