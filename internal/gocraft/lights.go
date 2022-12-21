package gocraft

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightType int32

const (
	LightTypeDirectional LightType = iota
	LightTypePoint
)

type Light struct {
	shader    rl.Shader
	lightType LightType
	position  rl.Vector3
	target    rl.Vector3
	color     rl.Color
	enabled   int32

	// shader locations
	enabledLoc int32
	typeLoc    int32
	posLoc     int32
	targetLoc  int32
	colorLoc   int32
}

const maxLightsCount = 4

var lightCount = 0

func NewLight(
	lightType LightType,
	position, target rl.Vector3,
	color rl.Color,
	shader rl.Shader) Light {

	light := Light{
		shader: shader,
	}

	if lightCount < maxLightsCount {
		light.enabled = 1
		light.lightType = lightType
		light.position = position
		light.target = target
		light.color = color

		light.enabledLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].enabled", lightCount))
		light.typeLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].type", lightCount))
		light.posLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].position", lightCount))
		light.targetLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].target", lightCount))
		light.colorLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].color", lightCount))
		light.UpdateValues()

		lightCount++
	}

	return light
}

// Send light properties to shader
func (lt *Light) UpdateValues() {
	setShaderInt32(lt.shader, lt.enabledLoc, lt.enabled)
	setShaderInt32(lt.shader, lt.enabledLoc, int32(lt.lightType))
	setShaderVec3(lt.shader, lt.posLoc, lt.position)
	setShaderVec3(lt.shader, lt.targetLoc, lt.target)
	setShaderColor(lt.shader, lt.colorLoc, lt.color)
}
