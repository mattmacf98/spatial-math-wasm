// This serves as the dummy version of what spatial math could look like if we take out the dependency on the rdk utils (which is not WASM compilable)
package main

import (
	"fmt"

	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
	"gonum.org/v1/gonum/num/dualquat"
	"gonum.org/v1/gonum/num/quat"

	"github.com/golang/geo/r3"
)

type Pose interface {
	Point() r3.Vector
	Orientation() Orientation
}

type Quaternion quat.Number
type DualQuaternion struct {
	dualquat.Number
}

func (q *DualQuaternion) Orientation() Orientation {
	return (*Quaternion)(&q.Real)
}

func (q *DualQuaternion) Point() r3.Vector {
	tQuat := dualquat.Mul(q.Number, dualquat.Conj(q.Number)).Dual
	return r3.Vector{tQuat.Imag, tQuat.Jmag, tQuat.Kmag}
}

// newDualQuaternion returns a pointer to a new dualQuaternion object whose Quaternion is an identity Quaternion.
// Since the real part of a qual quaternion should be a unit quaternion, not all zeroes, this should be used
// instead of &dualQuaternion{}.
func newDualQuaternion() *DualQuaternion {
	return &DualQuaternion{dualquat.Number{
		Real: quat.Number{Real: 1},
		Dual: quat.Number{},
	}}
}
func NewZeroPose() Pose {
	return newDualQuaternion()
}

type Path []referenceframe.FrameSystemPoses
type Trajectory []referenceframe.FrameSystemInputs
type SimplePlan struct {
	path Path
	traj Trajectory
}

func (plan *SimplePlan) Path() Path {
	return plan.path
}

func TrajectoryFromLinearInputs(inps []*referenceframe.LinearInputs) Trajectory {
	ret := make(Trajectory, len(inps))
	for idx, inp := range inps {
		ret[idx] = inp.ToFrameSystemInputs()
	}

	return ret
}

func NewSimplePlanFromTrajectory(
	trajAsInputs []*referenceframe.LinearInputs, fs *referenceframe.FrameSystem,
) (*SimplePlan, error) {
	path := Path{}
	for _, inputNode := range trajAsInputs {
		poseMap := make(map[string]*referenceframe.PoseInFrame)
		for frame := range inputNode.Keys() {
			tf, err := fs.Transform(inputNode, referenceframe.NewPoseInFrame(frame, spatialmath.NewZeroPose()), referenceframe.World)
			if err != nil {
				return nil, err
			}
			pose, ok := tf.(*referenceframe.PoseInFrame)
			if !ok {
				return nil, fmt.Errorf("pose not transformable")
			}
			poseMap[frame] = pose
		}
		path = append(path, poseMap)
	}

	return &SimplePlan{path: path, traj: TrajectoryFromLinearInputs(trajAsInputs)}, nil
}

type Orientation interface {
}

func (path Path) GetFramePoses(frameName string) ([]Pose, error) {
	poses := []Pose{}
	for _, step := range path {
		poseInFrame, ok := step[frameName]
		if !ok {
			return nil, fmt.Errorf("frame named %s not found in path", frameName)
		}
		pose := poseInFrame.Pose().(Pose)
		poses = append(poses, pose)
	}
	return poses, nil
}

func GetPosesFromTrajectory(
	fs *referenceframe.FrameSystem,
	trajectory Trajectory,
	frameName string,
) ([]Pose, error) {
	linearInputs := []*referenceframe.LinearInputs{}
	for _, trajInput := range trajectory {
		linearInputs = append(linearInputs, trajInput.ToLinearInputs())
	}
	sp, err := NewSimplePlanFromTrajectory(linearInputs, fs)
	if err != nil {
		return nil, fmt.Errorf("failed to create simple plan from trajectory: %w", err)
	}

	poses, err := sp.Path().GetFramePoses(frameName)
	if err != nil {
		return nil, fmt.Errorf("failed to get frame poses: %w", err)
	}

	return poses, nil
}
