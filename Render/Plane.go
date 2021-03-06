package Render

import (
	of "OctaForceEngineGo"
	"github.com/go-gl/mathgl/mgl32"
)

var plane []int
var baseColor mgl32.Vec3

const (
	size  = 100
	scale = 10
)

func getPointInPlane(x int, y int) int {
	return plane[x*size*2+y]
}
func setPointInPlane(x int, y int, id int) {
	plane[x*size*2+y] = id
}
func setUpPlane() {
	plane = make([]int, size*size*4)
	baseColor = mgl32.Vec3{0.5, 0.5, 0.5}

	point := of.CreateEntity()
	mesh := of.Mesh{}
	mesh.LoadOBJ(absPath+"/mesh/LowPolySphere.obj", false)
	mesh.Material = of.Material{DiffuseColor: baseColor}
	of.AddComponent(point, of.ComponentMesh, mesh)

	pointTransform := of.GetComponent(point, of.ComponentTransform).(of.Transform)
	pointTransform.Position = mgl32.Vec3{float32(-size) * scale, 0, float32(-size) * scale}
	of.SetComponent(point, of.ComponentTransform, pointTransform)

	setPointInPlane(0, 0, point)

	for x := 0; x < size*2; x++ {
		for y := 0; y < size*2; y++ {
			if x == 0 && y == 0 {
				continue
			}
			particle := of.CreateEntity()

			particleTransform := of.Transform{}
			particleTransform.Position = mgl32.Vec3{float32(x-size) * scale, 0, float32(y-size) * scale}
			of.AddComponent(particle, of.ComponentTransform, particleTransform)

			meshInstant := of.MeshInstant{}
			meshInstant.MeshEntity = plane[0]
			meshInstant.Material = of.Material{DiffuseColor: baseColor}
			of.AddComponent(particle, of.ComponentMeshInstant, meshInstant)

			setPointInPlane(x, y, particle)
		}
	}
}

var gravityMultiplier float32 = 100000.0
var cancelGravitation = false
var waveSpeed float32 = 1.0

func updatePlane() {

	for _, point := range plane {
		transform := of.GetComponent(point, of.ComponentTransform).(of.Transform)
		pos := transform.Position
		pos[1] = 0

		distance := mgl32.Vec3{0, 0, 0}.Sub(mgl32.Vec3{pos[0], 0, pos[2]}).Len()
		if cancelGravitation {
			pos[1] += 2.0 / distance * gravityMultiplier
		}
		timeOfWave := frame - int(distance/waveSpeed)

		posG0 := getPosOfGravityObject(timeOfWave, 0)
		distanceG0 := posG0.Sub(mgl32.Vec3{pos[0], 0, pos[2]}).Len()
		pos[1] -= 1 / distanceG0 * gravityMultiplier

		posG1 := getPosOfGravityObject(timeOfWave, 1)
		distanceG1 := posG1.Sub(mgl32.Vec3{pos[0], 0, pos[2]}).Len()
		pos[1] -= 1 / distanceG1 * gravityMultiplier

		transform.Position = pos
		of.SetComponent(point, of.ComponentTransform, transform)
	}

}
