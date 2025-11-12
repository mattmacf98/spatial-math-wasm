module spatial-math-wasm

go 1.25.3

require (
	github.com/go-gl/mathgl v1.0.0
	gonum.org/v1/gonum v0.16.0
)

require golang.org/x/image v0.25.0 // indirect

replace go.viam.com/utils => ./local/go.viam.com/utils
