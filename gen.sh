#!/bin/bash

base_path="$(dirname $0)"

idl_dir_path="${base_path}/idl"
idl_file_name="simple_tiktok.thrift"
idl_path="${idl_dir_path}/${idl_file_name}"
mod_name="simple-tiktok"
out_dir="${base_path}"

usage() {
cat << EOF
Usage
    gen.sh [OPTIONS ...]
    This is script used generate gateway apiservice code
OPTIONS
    -h, --help              display this help and exit
    -n, --new               first generate code
    -u, --update            update code
EXAMPLE
    bash gen.sh -n
EOF
exit 0
}

case "$1" in
    -h|--help)
        usage
        ;;
    -n|--new)
        # echo ${mod_name} ${idl_path}
        # echo "hz new -mod ${mod_name} -idl ${idl_path} -out_dir ${out_dir}"
        hz new -mod ${mod_name} -idl ${idl_path} -out_dir ${out_dir}
        ;;
    -u|--update)
        # echo "hz update -idl ${idl_path} -out_dir ${out_dir}"
        hz update -idl ${idl_path} -out_dir ${out_dir}
        ;;
    *)
        usage
        ;;
esac
