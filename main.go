// for functionality we want to turn into WASM, we need this entry point file to map the values from JS , through go and back to JS (we could have a single one for all rdk we want to WASMize or split out into one per repo (or these could be in the repos themselves)
package main

import (
	"syscall/js"

	"gonum.org/v1/gonum/num/quat"
)

func quatToOv(this js.Value, args []js.Value) any {
	w := args[0].Float()
	x := args[1].Float()
	y := args[2].Float()
	z := args[3].Float()

	q := quat.Number{Real: w, Imag: x, Jmag: y, Kmag: z}
	ov := QuatToOV(q)

	return js.ValueOf(map[string]any{
		"theta": ov.Theta,
		"ox":    ov.OX,
		"oy":    ov.OY,
		"oz":    ov.OZ,
	})
}

func main() {
	js.Global().Set("quatToOv", js.FuncOf(quatToOv))
	select {}
}
