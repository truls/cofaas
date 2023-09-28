SRC_DEPS = component.go ../../wit/chained-service.wit

ADAPTER = ../../deps/wasmtime/target/wasm32-unknown-unknown/release/wasi_preview1_component_adapter.wasm

all:main-component.wasm

main.wasm:${SRC_DEPS}
	go generate
	tinygo build -target=wasi -o main.wasm component.go

main.embed.wasm:main.wasm
	wasm-tools component embed --world consumer-interface ../..//wit main.wasm -o main.embed.wasm

main-component.wasm:main.embed.wasm ${ADAPTER}
	wasm-tools component new main.embed.wasm --adapt wasi_snapshot_preview1=${ADAPTER} -o main-component.wasm
	wasm-tools validate main-component.wasm --features component-model

clean:
	rm -rf gen/ main.wasm main.embed.wasm main-component.wasm

.PHONY:clean
