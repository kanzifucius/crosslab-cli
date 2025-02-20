#!/bin/sh

set -e

if ! command -v gh >/dev/null 2>&1; then
    echo "gh is not installed. Visit https://github.com/cli/cli#installation for installation instructions."
    exit 1
fi

get_arch() {
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" )
            echo "amd64"
        ;;
        "i386" | "i486" | "i586")
            echo "386"
        ;;
        "aarch64" | "arm64" | "arm")
            echo "arm64"
        ;;
        "mips64el")
            echo "mips64el"
        ;;
        "mips64")
            echo "mips64"
        ;;
        "mips")
            echo "mips"
        ;;
        *)
            echo ${NIL}
        ;;
    esac
}

get_os(){
    # darwin: Darwin
    echo $(uname -s | awk '{print tolower($0)}')
}

os=$(get_os)
arch=$(get_arch)

gh release download --repo opus2-platform/platform-internal --pattern "platform_${os}_${arch}" -O platform --clobber
chmod 0755 ./platform
 