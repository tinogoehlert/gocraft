package gocraft

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	fastnoise "github.com/tinogoehlert/gocraft/pkg/go-fastnoiselite"
)

type Block struct {
	position    rl.Vector3
	heightValue float32
	heightNoise []float32
}

func normalizef(in float32) float32 {
	out := (in - -1) / (1 - -1) * (0 - 32)
	out = float32(math.Abs(float64(out)))

	return float32(math.Round(float64(out)))
}

type Chunk struct {
	Center     rl.Vector2
	BoudingBox rl.BoundingBox
	Width      float32
	Height     float32
	Lenght     float32
	blocks     []*Block
	debugColor rl.Color
}

type blockInstances struct {
	transforms []rl.Matrix
	material   rl.Material
}

func (c *Chunk) BuildMesh(blockMesh rl.Model, mat []rl.Material) {
	count := 0
	var (
		dirtBlocks = blockInstances{
			transforms: make([]rl.Matrix, 0),
			material:   mat[0],
		}
		grasBlocks = blockInstances{
			transforms: make([]rl.Matrix, 0),
			material:   mat[1],
		}
		snowBlocks = blockInstances{
			transforms: make([]rl.Matrix, 0),
			material:   mat[2],
		}
	)

	for _, block := range c.blocks {
		t := rl.MatrixTranslate(
			block.position.X,
			0,
			block.position.Z,
		)
		dirtBlocks.transforms = append(dirtBlocks.transforms, t)
		for i := 1; i <= int(block.heightValue); i++ {
			if block.heightNoise[i] < 20 {
				continue
			}
			block.position.Y = float32(i)
			t := rl.MatrixTranslate(
				block.position.X,
				block.position.Y,
				block.position.Z,
			)
			count++
			if i == int(block.heightValue) {
				if i > 25 {
					snowBlocks.transforms = append(snowBlocks.transforms, t)
				} else {
					grasBlocks.transforms = append(grasBlocks.transforms, t)
				}
			} else {
				dirtBlocks.transforms = append(dirtBlocks.transforms, t)
			}
		}
	}

	rl.DrawMeshInstanced(
		*blockMesh.Meshes,
		dirtBlocks.material,
		dirtBlocks.transforms,
		len(dirtBlocks.transforms),
	)

	if len(grasBlocks.transforms) > 0 {
		rl.DrawMeshInstanced(
			*blockMesh.Meshes,
			grasBlocks.material,
			grasBlocks.transforms,
			len(grasBlocks.transforms),
		)
	}

	if len(snowBlocks.transforms) > 0 {
		rl.DrawMeshInstanced(
			*blockMesh.Meshes,
			snowBlocks.material,
			snowBlocks.transforms,
			len(snowBlocks.transforms),
		)
	}
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
	cm.terrainNoise.SetType(fastnoise.FNL_NOISE_OPENSIMPLEX2S)
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
		dist := rl.Vector2Distance(oldChunk.Center, rl.NewVector2(-pos.X, -pos.Z))
		if dist < cm.length/2 {
			rl.DrawText(fmt.Sprintf("distance: %d", int(dist)), 10, 60, 16, oldChunk.debugColor)
			rl.DrawText(fmt.Sprintf("current chunk: %v", oldChunk.Center), 10, 20, 16, oldChunk.debugColor)
			break
		}
	}
}

func (cm *ChunkManager) GetChunks(pos rl.Vector3, source *Chunk) []*Chunk {
	currentChunk = source
	for _, oldChunk := range cm.chunkMap {
		dist := rl.Vector2Distance(oldChunk.Center, rl.NewVector2(-pos.X, -pos.Z))
		if dist < cm.length/2 {
			currentChunk = oldChunk
			break
		}
	}

	var (
		chunks = []*Chunk{}
		west   = rl.NewVector2(
			currentChunk.Center.X+(cm.width),
			currentChunk.Center.Y,
		)
		east = rl.NewVector2(
			(currentChunk.Center.X - (cm.width)),
			currentChunk.Center.Y,
		)
		eastwest = rl.NewVector2(
			(currentChunk.Center.X - (cm.width)),
			(currentChunk.Center.Y - (cm.width)),
		)
		southwest = rl.NewVector2(
			(currentChunk.Center.X + (cm.width)),
			(currentChunk.Center.Y + (cm.width)),
		)
		northeast = rl.NewVector2(
			(currentChunk.Center.X - (cm.width)),
			(currentChunk.Center.Y + (cm.width)),
		)
		southeast = rl.NewVector2(
			(currentChunk.Center.X + (cm.width)),
			(currentChunk.Center.Y - (cm.width)),
		)
		north = rl.NewVector2(
			currentChunk.Center.X,
			currentChunk.Center.Y+(cm.height),
		)
		south = rl.NewVector2(
			currentChunk.Center.X,
			(currentChunk.Center.Y - (cm.height)),
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
	fmt.Println("generate new chunk: ", pos)
	var (
		halfWidth  = cm.width / 2
		halfLenght = cm.length / 2
		startX     = pos.X + halfWidth
		startY     = pos.Y + halfLenght
	)

	chunk := &Chunk{
		Width:  cm.width,
		Lenght: cm.length,
		Height: cm.height,
		Center: pos,
		blocks: make([]*Block, 0, int(cm.width*cm.length)),
		BoudingBox: rl.NewBoundingBox(
			rl.NewVector3(
				pos.X-halfWidth,
				0,
				pos.Y-halfLenght,
			),
			rl.NewVector3(
				pos.X+halfWidth,
				cm.height,
				pos.Y+halfLenght,
			),
		),
	}

	var l float32
	for l = 0; l < cm.length; l++ {
		var w float32
		for w = 0; w < cm.width; w++ {
			heightData := cm.terrainNoise.GetNoise2D(w-startX, l-startY)
			b := &Block{
				heightValue: normalizef(heightData),
				position:    rl.NewVector3(w-startX, 0, l-startY),
			}
			for i := 0; i <= int(b.heightValue); i++ {
				n := cm.terrainNoise.GetNoise3D(w-startX, float32(i), l-startY)
				b.heightNoise = append(b.heightNoise, normalizef(n))
			}
			chunk.blocks = append(chunk.blocks, b)
		}
	}

	cm.chunkMap[pos] = chunk
	return chunk
}
