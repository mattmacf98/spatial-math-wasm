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
	// Take the frameSystem as a JSON string and unmarshal into FrameSystem
	frameSystemJSON := args[0].String()
	var frameSystem referenceframe.FrameSystem
	err := json.Unmarshal([]byte(frameSystemJSON), &frameSystem)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal FrameSystem: %v", err))
	}
	trajectoryJSON := args[1].String()
	var trajectory motionplan.Trajectory
	err = json.Unmarshal([]byte(trajectoryJSON), &trajectory)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Trajectory: %v", err))
	}
	frameName := args[2].String()

	poses, err := GetPosesFromTrajectory(&frameSystem, trajectory, frameName)
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
