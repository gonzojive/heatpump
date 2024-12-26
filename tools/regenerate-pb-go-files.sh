#!/bin/bash

set -euxo pipefail

# cd to the root of the git repo
cd "$(git rev-parse --show-toplevel)"

# There should be one command for each result of
# bazel query 'kind("go_proto_library", //...)'

bazel run //proto/logs:write_generated_go_files
bazel run //proto/chiltrix:write_generated_go_files
bazel run //proto/fancoil:write_generated_go_files
