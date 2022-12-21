package gocraft

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ErrBlockManagerNoSpace = errors.New("not enough space to add block")
)

type BlockType int

const (
	BlockTypeDirt  BlockType = 0
	BlockTypeGras  BlockType = 1
	BlockTypeSnow  BlockType = 2
	BlockTypeRock  BlockType = 3
	BlockTypeGroud BlockType = 4
)

var blockMaterials = map[BlockType]rl.Material{}

func addBlockMaterial(t BlockType, fn string, shader rl.Shader) {
	bm := rl.LoadMaterialDefault()
	bm.Shader = shader
	bm.Maps.Texture = rl.LoadTexture(fn)
	blockMaterials[t] = bm
}

type Block struct {
	position  rl.Vector3
	blockType BlockType
	enabled   bool
	carved    bool
}
