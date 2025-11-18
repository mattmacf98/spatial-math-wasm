// for functionality we want to turn into WASM, we need this entry point file to map the values from JS , through go and back to JS (we could have a single one for all rdk we want to WASMize or split out into one per repo (or these could be in the repos themselves)
package main

import (
	"encoding/json"
	"fmt"
	"spatial-math-wasm/motionplan"
	"spatial-math-wasm/referenceframe"
	"spatial-math-wasm/spatialmath"
	"syscall/js"
)

func GetPosesFromTrajectory(
	fs *referenceframe.FrameSystem,
	trajectory motionplan.Trajectory,
	frameName string,
) ([]spatialmath.Pose, error) {
	linearInputs := []*referenceframe.LinearInputs{}
	for _, trajInput := range trajectory {
		linearInputs = append(linearInputs, trajInput.ToLinearInputs())
	}
	sp, err := motionplan.NewSimplePlanFromTrajectory(linearInputs, fs)
	if err != nil {
		return nil, fmt.Errorf("failed to create simple plan from trajectory: %w", err)
	}

	poses, err := sp.Path().GetFramePoses(frameName)
	if err != nil {
		return nil, fmt.Errorf("failed to get frame poses: %w", err)
	}

	return poses, nil
}

func getPosesFromTrajectory(this js.Value, args []js.Value) any {
	frameSystemString := args[0].String()
	trajectoryString := args[1].String()
	frameName := args[2].String()

	fmt.Printf("got here 1\n")

	var frameSystem referenceframe.FrameSystem
	err := json.Unmarshal([]byte(frameSystemString), &frameSystem)
	if err != nil {
		return fmt.Errorf("failed to unmarshal frame system: %w", err)
	}
	fmt.Printf("got here 2\n")

	var trajectory motionplan.Trajectory
	err = json.Unmarshal([]byte(trajectoryString), &trajectory)
	if err != nil {
		return fmt.Errorf("failed to unmarshal trajectory: %w", err)
	}
	fmt.Printf("got here 3\n")

	fmt.Printf("FRAME SYSTEM: %v\n", frameSystem)
	fmt.Printf("TRAJECTORY: %v\n", trajectory)
	fmt.Printf("FRAME NAME: %s\n", frameName)

	poses, err := GetPosesFromTrajectory(&frameSystem, trajectory, frameName)
	if err != nil {
		return js.ValueOf(map[string]any{
			"error": err.Error(),
		})
	}
	fmt.Printf("got here 4\n")
	jsonPoses, _ := json.Marshal(poses)
	fmt.Printf("JSON POSES: %s\n", string(jsonPoses))
	return js.Global().Get("JSON").Call("parse", string(jsonPoses))
}

func main() {
	js.Global().Set("getPosesFromTrajectory", js.FuncOf(getPosesFromTrajectory))
	select {}
}
