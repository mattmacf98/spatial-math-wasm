// for functionality we want to turn into WASM, we need this entry point file to map the values from JS , through go and back to JS (we could have a single one for all rdk we want to WASMize or split out into one per repo (or these could be in the repos themselves)
package main

import (
	"syscall/js"

	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

func getPosesFromTrajectory(this js.Value, args []js.Value) any {
	frameSystem := args[0].Interface().(*referenceframe.FrameSystem)
	trajectory := args[1].Interface().(*motionplan.Trajectory)
	frameName := args[2].String()

	poses, err := GetPosesFromTrajectory(frameSystem, trajectory, frameName)
	if err != nil {
		return js.ValueOf(map[string]any{
			"error": err.Error(),
		})
	}

	return js.ValueOf(poses)
}

func main() {
	js.Global().Set("getPosesFromTrajectory", js.FuncOf(getPosesFromTrajectory))
	select {}
}
