1. GOOS=js GOARCH=wasm go build -o main.wasm
2. cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
3. we need to make a wrapper main to export all the functions
4. then we can load the wasm instance up and it will auto-inject our functions

# Issues

1. the spatial math repo has an import to viam utils which cannot be compiled to WASM

```
../../../go/pkg/mod/go.viam.com/utils@v0.1.176/runtime.go:38:24: undefined: syscall.SIGPIPE
 ../../../go/pkg/mod/go.viam.com/utils@v0.1.176/runtime_posix.go:12:33: undefined: syscall.SIGUSR1
```

we could probably refactor a lot of the functions we need out but might/ mark files as not to be compiled for WASM in rdk but will be a lot

# TODO

1. build out an example svelte kit app that uses this spatial math and properly gets the decleration
2. assess what actual functions we need to call for our P0 scope
