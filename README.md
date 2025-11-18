# Spatial Math WASM

## About

This is a POC of how we can use WASM to compile some math functions in RDK and use them client-side in apps.

The important processes are:

1. Creating a wrapper Go entry in [`main.go`](./main.go) to pass JS → Go → JS
2. A [`wasm.ts`](./spatial-math-wasm/src/lib/wasm.ts) file which handles initializing the WASM instance and exposes wrappers around the WASM execution (to ensure WASM is ready, the function exists in WASM, and the values won't cause the function to panic)
3. In [`+page.svelte`](./spatial-math-wasm/src/routes/+page.svelte) we can see how we are then able to call the function as if it were a function declared in TypeScript

## Build Instructions

1. Build the WASM binary:

   ```bash
   GOOS=js GOARCH=wasm go build -o ./spatial-math-wasm-web/static/main.wasm
   ```

2. Copy the WASM JavaScript runtime:

   ```bash
   cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
   ```

3. Create a wrapper main to export all the functions

4. Load the WASM instance and it will auto-inject our functions

## Issues

1. The spatial math repo has an import to viam utils which cannot be compiled to WASM:

   ```
   ../../../go/pkg/mod/go.viam.com/utils@v0.1.176/runtime.go:38:24: undefined: syscall.SIGPIPE
   ../../../go/pkg/mod/go.viam.com/utils@v0.1.176/runtime_posix.go:12:33: undefined: syscall.SIGUSR1
   ```

   We could probably refactor a lot of the functions we need out or mark files as not to be compiled for WASM in RDK, but this will require significant work.

   I created a minimum working version, a lot of the work was around removing references to `go.viam.com/utils` for simple math conversions and logging. There was also a decent amount of test removing since they used some file loading things in `go.viam.com/utils` as well.

## TODO

- [ ] Try to make `useWasm` a hook

## Instructions

1. get a set of `frame_system.json` and `plans.json` from running the sanding cli (or use the ones in the examples directory)
2. open the console (for this poc I just dumped the output poses in logs)
3. click the `Get Poses` button and see the first trajectory + frame system gets converted into a list of poses (to change what # trajectory gets converted change the indices into trajectories and trajectory in the `testWasm` function in `+page.svelte`)
