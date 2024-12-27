#!/bin/bash

set -euxo pipefail

# cd to the root of the git repo
cd "$(git rev-parse --show-toplevel)"

# There should be one command for each result of
# bazel query 'kind("go_proto_library", //...)'

bazel run //proto/chiltrix:chiltrix_go_proto_write_source_files
bazel run //proto/command_queue:command_queue_go_proto_write_source_files
bazel run //proto/controller:controller_go_proto_write_source_files
bazel run //proto/fancoil:fancoil_go_proto_write_source_files
bazel run //proto/logs:logs_go_proto_write_source_files
