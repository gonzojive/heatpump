#!/bin/bash

set -e # fail if anything exits non zero

bazel run //cmd/authservice:push-image
bazel run //cmd/cloud/http-endpoint:push_image
bazel run //cmd/queueserver:push-image
bazel run //cmd/reverse-proxy:push-image
bazel run //cmd/stateservice:push-image
