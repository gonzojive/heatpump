"""A replacement for go_proto_library that outputs generated files to the
source directory.
"""

load(":write_go_generated_source_files.bzl", "write_go_generated_source_files")
load("@io_bazel_rules_go//proto:def.bzl", orig_go_proto_library = "go_proto_library")

# buildifier: disable=function-docstring-args
def go_proto_library(name, output_files, **kwargs):
    """Wrapper around go_proto_library that outputs generated sources to the source
    directory.
    """
    orig_go_proto_library(name = name, **kwargs)
    write_go_generated_source_files(
        name + "_write_source_files",
        output_files = output_files,
        src = name,
    )