#!/bin/bash

set -euxo pipefail

# cd to the root of the git repo
cd "$(git rev-parse --show-toplevel)"

bazel build //cmd/chiltrix:chiltrix_release

TAR="bazel-bin/cmd/chiltrix/chiltrix_release.tar"

# Transfer and extract the tar archive to the remote directory
cat "$TAR" | ssh waterpi "tar -xf - -C /home/pi/bin"

# Execute the command on the remote host with Ctrl-C forwarding
ssh -t waterpi "/home/pi/bin/chiltrix start-readwrite-service --alsologtostderr --grpc-port 8084"
