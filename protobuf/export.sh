#!/bin/sh

# ************************* 解析参数 *************************
usage="usage: $0 [-w workspace][-w] [-s code_sign_identity] [-h help]"

current_shell_path=$(dirname $0)
current_shell_path=${current_shell_path/\./$(pwd)}
export_path=$current_shell_path/../model/gen

cd $current_shell_path/


rm -rf ${export_path}
mkdir -pv ${export_path}
echo "❤️  protobuf 将导出到 ${export_path}"
echo 'protobuf 开始生成 ...'

$current_shell_path/protoc/bin/protoc *.proto --go_out=$export_path

echo " ❤️  protobuf 导出完成"

exit 0