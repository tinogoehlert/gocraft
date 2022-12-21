package gocraft

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	fastnoise "github.com/tinogoehlert/gocraft/pkg/go-fastnoiselite"
)

func normalizef(in float32) float32 {
	out := (in - -1) / (1 - -1) * (0 - 48)
	out = float32(math.Abs(float64(out)))

	return float32(math.Round(float64(out)))
}

type ChunkManager struct {
	terrainNoise *fastnoise.NoiseState
	chunkMap     map[rl.Vector2]*Chunk
	width        float32
	height       float32
	length       float32
}

func NewChunkManager(size float32) *ChunkManager {
	cm := &ChunkManager{
		terrainNoise: fastnoise.NewDefaultNoise(),
		chunkMap:     make(map[rl.Vector2]*Chunk),
		width:        size,
		height:       size,
		length:       size,
	}

	cm.terrainNoise.SetFractal(fastnoise.FNL_FRACTAL_RIDGED)
	cm.terrainNoise.SetType(fastnoise.FNL_NOISE_PERLIN)
	cm.terrainNoise.SetFrequency(0.01)
	cm.terrainNoise.SetOctaves(4)
	return cm
}

func (cm *ChunkManager) GetChunk(pos rl.Vector2, color rl.Color) *Chunk {
	if chunk, ok := cm.chunkMap[pos]; ok {
		return chunk
	}

	c := cm.generateNewChunk(pos)
	c.debugColor = color
	return c
}

var currentChunk *Chunk

func (cm ChunkManager) DebugChunks(pos rl.Vector3, source *Chunk) {
	rl.DrawText(fmt.Sprintf("position: %v", rl.NewVector2(pos.X, pos.Z)), 10, 40, 16, rl.Yellow)
	for _, oldChunk := range cm.chunkMap {
		dist := rl.Vector2Distance(oldChunk.center, rl.NewVector2(-pos.X, -pos.Z))
		if dist < cm.length/2 {
			rl.DrawText(fmt.Sprintf("distance: %d", int(dist)), 10, 60, 16, oldChunk.debugColor)
			rl.DrawText(fmt.Sprintf("current chunk: %v", oldChunk.center), 10, 20, 16, oldChunk.debugColor)
			break
		}
	}
}

func (cm *ChunkManager) GetChunks(pos rl.Vector3, source *Chunk) []*Chunk {
	currentChunk = source
	for _, oldChunk := range cm.chunkMap {
		dist := rl.Vector2Distance(oldChunk.center, rl.NewVector2(-pos.X, -pos.Z))
		if dist < cm.length/2 {
			currentChunk = oldChunk
			break
		}
	}

	var (
		chunks = []*Chunk{}
		west   = rl.NewVector2(
			currentChunk.center.X+(cm.width),
			currentChunk.center.Y,
		)
		east = rl.NewVector2(
			(currentChunk.center.X - (cm.width)),
			currentChunk.center.Y,
		)
		eastwest = rl.NewVector2(
			(currentChunk.center.X - (cm.width)),
			(currentChunk.center.Y - (cm.width)),
		)
		southwest = rl.NewVector2(
			(currentChunk.center.X + (cm.width)),
			(currentChunk.center.Y + (cm.width)),
		)
		northeast = rl.NewVector2(
			(currentChunk.center.X - (cm.width)),
			(currentChunk.center.Y + (cm.width)),
		)
		southeast = rl.NewVector2(
			(currentChunk.center.X + (cm.width)),
			(currentChunk.center.Y - (cm.width)),
		)
		north = rl.NewVector2(
			currentChunk.center.X,
			currentChunk.center.Y+(cm.height),
		)
		south = rl.NewVector2(
			currentChunk.center.X,
			(currentChunk.center.Y - (cm.height)),
		)
	)

	chunks = append(chunks,
		currentChunk,
		cm.GetChunk(west, rl.Green),
		cm.GetChunk(east, rl.Yellow),
		cm.GetChunk(east, rl.Yellow),
		cm.GetChunk(north, rl.Blue),
		cm.GetChunk(south, rl.Red),
		cm.GetChunk(eastwest, rl.Purple),
		cm.GetChunk(northeast, rl.Pink),
		cm.GetChunk(southeast, rl.Brown),
		cm.GetChunk(southwest, rl.Orange),
	)
	return chunks
}

func (cm *ChunkManager) generateNewChunk(pos rl.Vector2) *Chunk {
	var (
		halfWidth  = cm.width / 2
		halfLenght = cm.length / 2
		startX     = pos.X + halfWidth
		startY     = pos.Y + halfLenght
	)

	chunk := NewChunk(int(cm.width), int(cm.height), int(cm.length))
	chunk.center = pos

	var l float32
	for l = 0; l < cm.width; l++ {
		var w float32
		for w = 0; w < cm.height; w++ {
			var h float32
			cm.terrainNoise.SetFractal(fastnoise.FNL_FRACTAL_RIDGED)
			cm.terrainNoise.SetOctaves(4)
			terraXZ := normalizef(cm.terrainNoise.GetNoise2D(w-startX, l-startY))
			for h = 0; h < cm.height; h++ {
				if h < terraXZ || h == 0 {
					cm.terrainNoise.SetFractal(fastnoise.FNL_FRACTAL_PINGPONG)
					cm.terrainNoise.SetOctaves(6)
					terraXYZ := normalizef(cm.terrainNoise.GetNoise3D(w-startX, float32(h), l-startY))
					b := &Block{
						blockType: BlockTypeDirt,
						position:  rl.NewVector3(w-startX, h, l-startY),
					}
					if terraXYZ < 10 && h > 0 {
						b.carved = true
					}
					chunk.AddBlock(b, w, h, l)
				}
			}
		}
	}

	chunk.Generate()
	cm.chunkMap[pos] = chunk
	return chunk
}
