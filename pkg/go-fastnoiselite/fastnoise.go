package fastnoise

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #define FNL_IMPL
// #include "fastnoise.h"
import "C"

// Enums
type FNL_NOISE int
type FNL_ROTATION int
type FNL_FRACTAL int
type FNL_CELLULAR_DISTANCE int
type FNL_CELLULAR_RETURN_VALUE int
type FNL_DOMAIN_WARP int

const (
	FNL_NOISE_OPENSIMPLEX2  FNL_NOISE = 0
	FNL_NOISE_OPENSIMPLEX2S FNL_NOISE = 1
	FNL_NOISE_CELLULAR      FNL_NOISE = 2
	FNL_NOISE_PERLIN        FNL_NOISE = 3
	FNL_NOISE_VALUE_CUBIC   FNL_NOISE = 4
	FNL_NOISE_VALUE         FNL_NOISE = 5
)

const (
	FNL_ROTATION_NONE              FNL_ROTATION = 0
	FNL_ROTATION_IMPROVE_XY_PLANES FNL_ROTATION = 1
	FNL_ROTATION_IMPROVE_XZ_PLANES FNL_ROTATION = 2
)

const (
	FNL_FRACTAL_NONE                    FNL_FRACTAL = 0
	FNL_FRACTAL_FBM                     FNL_FRACTAL = 1
	FNL_FRACTAL_RIDGED                  FNL_FRACTAL = 2
	FNL_FRACTAL_PINGPONG                FNL_FRACTAL = 3
	FNL_FRACTAL_DOMAIN_WARP_PROGRESSIVE FNL_FRACTAL = 4
	FNL_FRACTAL_DOMAIN_WARP_INDEPENDENT FNL_FRACTAL = 5
)

const (
	FNL_CELLULAR_DISTANCE_EUCLIDEAN   FNL_CELLULAR_DISTANCE = 0
	FNL_CELLULAR_DISTANCE_EUCLIDEANSQ FNL_CELLULAR_DISTANCE = 1
	FNL_CELLULAR_DISTANCE_MANHATTAN   FNL_CELLULAR_DISTANCE = 2
	FNL_CELLULAR_DISTANCE_HYBRI       FNL_CELLULAR_DISTANCE = 3
)

const (
	FNL_CELLULAR_RETURN_VALUE_CELLVALUE    FNL_CELLULAR_RETURN_VALUE = 0
	FNL_CELLULAR_RETURN_VALUE_DISTANCE     FNL_CELLULAR_RETURN_VALUE = 1
	FNL_CELLULAR_RETURN_VALUE_DISTANCE2    FNL_CELLULAR_RETURN_VALUE = 2
	FNL_CELLULAR_RETURN_VALUE_DISTANCE2ADD FNL_CELLULAR_RETURN_VALUE = 3
	FNL_CELLULAR_RETURN_VALUE_DISTANCE2SUB FNL_CELLULAR_RETURN_VALUE = 4
	FNL_CELLULAR_RETURN_VALUE_DISTANCE2MUL FNL_CELLULAR_RETURN_VALUE = 5
	FNL_CELLULAR_RETURN_VALUE_DISTANCE2DIV FNL_CELLULAR_RETURN_VALUE = 6
)

const (
	FNL_DOMAIN_WARP_OPENSIMPLEX2         FNL_DOMAIN_WARP = 0
	FNL_DOMAIN_WARP_OPENSIMPLEX2_REDUCED FNL_DOMAIN_WARP = 1
	FNL_DOMAIN_WARP_BASICGRID            FNL_DOMAIN_WARP = 2
)

type NoiseState struct {
	cstate C.struct_fnl_state
	warpX  C.float
	warpY  C.float
}

func NewDefaultNoise() *NoiseState {
	return &NoiseState{
		cstate: C.fnlCreateState(),
	}
}

func (n *NoiseState) SetSeed(seed int) {
	n.cstate.seed = C.int(seed)
}

func (n *NoiseState) SetOctaves(octaves int) {
	n.cstate.octaves = C.int(octaves)
}

func (n *NoiseState) SetFractal(frac FNL_FRACTAL) {
	n.cstate.fractal_type = C.fnl_fractal_type(frac)
}

func (n *NoiseState) SetType(typ FNL_NOISE) {
	n.cstate.noise_type = C.fnl_noise_type(typ)
}

func (n *NoiseState) SetFrequency(frequency float32) {
	n.cstate.frequency = C.float(frequency)
}

func (n *NoiseState) SetGain(gain float32) {
	n.cstate.gain = C.float(gain)
}

func (n *NoiseState) SetPingPongStrength(strength float32) {
	n.cstate.ping_pong_strength = C.float(strength)
}

// 3D noise at given position using the state settings
func (n *NoiseState) GetNoise2D(x, y float32) float32 {
	return float32(C.fnlGetNoise2D(&n.cstate, C.float(x), C.float(y)))
}

// 3D noise at given position using the state settings
func (n *NoiseState) DomainWarp2D(x, y float32) float32 {
	n.warpX = C.float(x)
	n.warpY = C.float(y)
	C.fnlDomainWarp2D(&n.cstate, &n.warpX, &n.warpY)
	return float32(C.fnlGetNoise2D(&n.cstate, n.warpX, n.warpY))
}

// 2D noise at given position using the state settings
func (n *NoiseState) GetNoise3D(x, y, z float32) float32 {
	return float32(C.fnlGetNoise3D(&n.cstate, C.float(x), C.float(y), C.float(z)))
}
