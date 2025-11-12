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
	<-make(chan bool)
}
