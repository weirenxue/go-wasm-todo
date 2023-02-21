.PHONY: build
build: clean
	mkdir ./build
	GOARCH=wasm GOOS=js go build -o ./build/main.wasm
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ./build/

	
.PHONY: clean
clean:
	rm -rf ./build