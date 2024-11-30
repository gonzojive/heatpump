#!/bin/bash

set -euxo pipefail

# cd to the root of the git repo
cd "$(git rev-parse --show-toplevel)"

bazel build //cmd/cx34install:cx34control_release

TAR="bazel-bin/cmd/cx34install/cx34control_release.tar"

# Transfer and extract the tar archive to the remote directory
cat "$TAR" | ssh waterpi "tar -xf - -C /home/pi/bin" 
