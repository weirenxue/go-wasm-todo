# Golang WebAssembly TODO Web Application

This is a project to explore the Golang Web Assembly by building a TODO application.

## How to build?

1. Install Golang version 1.19 or higher.

1. Use `make` command

    ```sh
    $ make build
    rm -rf ./build
    mkdir ./build
    GOARCH=wasm GOOS=js go build -o ./build/main.wasm
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./build/
    ```

1. The `index.html` is the home page.

---

## You can find more details in the [Blog](https://weirenxue.github.io/2023/02/21/go-wasm/)
