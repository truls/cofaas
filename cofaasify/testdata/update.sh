#!/bin/bash

set -xeuo pipefail

my_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
base_dir=${my_dir}/../../

protocmd_dir="$base_dir/lang-support/go/protocmd/"
protocmd="${protocmd_dir}"protocmd
(cd "$protocmd_dir" && go build)

_out_dir=$(mktemp -ud)
go build -C "$protocmd_dir"
(cd "$protocmd_dir" && "$protocmd" -exportProto "$my_dir"/protos/helloworld.proto -importProto "$my_dir"/protos/prodcon.proto -outputDir "$_out_dir" -witPath "$my_dir"/../../wit -witWorld "producer-interface" -implPath "$my_dir"/functions/go/producer)
cp -r "$_out_dir"/* "$my_dir"/component/go/producer/.
rm -r "$_out_dir"

_out_dir=$(mktemp -ud)
(cd "$protocmd_dir" && "$protocmd" -exportProto "$my_dir"/protos/prodcon.proto -outputDir "$_out_dir" -witPath "$my_dir"/../../wit -witWorld "consumer-interface" -implPath "$my_dir"/functions/go/consumer)
cp -r "$_out_dir"/* "$my_dir"/component/go/consumer/.
rm -r "$_out_dir"

rm -f "$my_dir"/component/go/producer/component/main.wasm
make -C "$my_dir"/component/go/producer/component

rm -f "$my_dir"/component/go/consumer/component/main.wasm
make -C "$my_dir"/component/go/consumer/component

(cd "$base_dir"/wrapper; cargo build --release)

make -C "$my_dir"/component/go
make -C "$my_dir"/functions/go

make -C "$my_dir"/component/rust
make -C "$my_dir"/functions/rust
