#!/bin/bash

set -euxo pipefail

# cd to the root of the git repo
cd "$(git rev-parse --show-toplevel)"

bazel build //cmd/cx34install:cx34control_release

TAR="bazel-bin/cmd/cx34install/cx34control_release.tar"

# Transfer and extract the tar archive to the remote directory
cat "$TAR" | ssh waterpi "tar -xf - -C /home/pi/bin" 

# Execute the command on the remote host with Ctrl-C forwarding
ssh -t waterpi "/home/pi/bin/cx34control --alsologtostderr --grpc-port 8084 --print-state-interval 10m"
