#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf class
javac --release 17 -g -parameters -d class $(find src -name '*.java')
