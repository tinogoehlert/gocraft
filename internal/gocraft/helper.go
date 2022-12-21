package gocraft

import (
	"math"
	"reflect"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func setShaderInt32(shader rl.Shader, locIndex int32, value int32) {
	// not pretty -_-
	// need nicer api
	sh := &reflect.SliceHeader{
		Len: 4,
		Cap: 4,
	}
	sh.Data = uintptr(unsafe.Pointer(&value))
	rl.SetShaderValue(shader, value, *(*[]float32)(unsafe.Pointer(sh)), rl.ShaderUniformInt)
}

func setShaderVec3(shader rl.Shader, locIndex int32, vec rl.Vector3) {
	rl.SetShaderValue(shader, locIndex, []float32{vec.X, vec.Y, vec.Z}, rl.ShaderUniformVec3)
}

func setShaderColor(shader rl.Shader, locIndex int32, col rl.Color) {
	rl.SetShaderValue(shader, locIndex, []float32{
		float32(col.R) / 255,
		float32(col.G) / 255,
		float32(col.B) / 255,
		float32(col.A) / 255,
	}, rl.ShaderUniformVec4)
}

func loglevelFromString(level string) rl.TraceLogLevel {
	switch level {
	case "all":
		return rl.LogAll
	case "trace":
		return rl.LogTrace
	case "debug":
		return rl.LogDebug
	case "info":
		return rl.LogInfo
	case "warn":
		return rl.LogWarning
	case "error":
		return rl.LogError
	case "fatal":
		return rl.LogFatal
	}

	return rl.LogError
}

func sin32(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func cos32(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func radians(deg float32) float32 {
	return deg * (math.Pi / 180)
}

func degrees(rad float32) float32 {
	return rad * (180 / math.Pi)
}
