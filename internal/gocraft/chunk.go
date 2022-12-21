package gocraft

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chunk struct {
	blockList  [][][]*Block
	instances  map[BlockType][]rl.Matrix
	center     rl.Vector2
	debugColor rl.Color
}

func NewChunk(width, height, lenght int) *Chunk {
	bm := &Chunk{}
	bm.blockList = make([][][]*Block, width)
	for i := 0; i < width; i++ {
		bm.blockList[i] = make([][]*Block, height)
		for j := 0; j < lenght; j++ {
			bm.blockList[i][j] = make([]*Block, lenght)
		}
	}
	bm.instances = make(map[BlockType][]rl.Matrix)
	return bm
}

func (bm *Chunk) GetBlock(x, y, z int) *Block {
	var (
		ix = int(x)
		iy = int(y)
		iz = int(z)
	)

	if err := bm.checkBlockPosition(ix, iy, iz); err != nil {
		return nil
	}

	return bm.blockList[ix][iy][iz]
}

func (bm *Chunk) HasBlock(x, y, z int) bool {
	b := bm.GetBlock(x, y, z)
	return b != nil && !b.carved
}

func (bm *Chunk) IsEnabled(x, y, z int) bool {
	b := bm.GetBlock(x, y, z)
	return b != nil && b.enabled && !b.carved
}

func (bm *Chunk) IsCarved(x, y, z int) bool {
	b := bm.GetBlock(x, y, z)
	return b != nil && b.carved
}

func (bm *Chunk) IsSurrounded(x, y, z int) bool {
	front := bm.HasBlock(x+1, y, z) && bm.HasBlock(x, y+1, z) && bm.HasBlock(x, y, z+1) && bm.HasBlock(x+1, y+1, z)
	back := bm.HasBlock(x-1, y, z) && bm.HasBlock(x, y-1, z) && bm.HasBlock(x, y, z-1) && bm.HasBlock(x-1, y-1, z)
	return front && back
}

func (bm *Chunk) IsSurroundedByCarved(x, y, z int) bool {
	front := bm.HasBlock(x+1, y, z) || bm.HasBlock(x, y, z+1)
	back := bm.HasBlock(x-1, y, z) || bm.HasBlock(x, y-1, z) || bm.HasBlock(x, y, z-1) || bm.HasBlock(x-1, y-1, z)
	return front && back
}

func (bm *Chunk) AddBlock(block *Block, x, y, z float32) error {
	var (
		ix = int(x)
		iy = int(y)
		iz = int(z)
	)

	if err := bm.checkBlockPosition(ix, iy, iz); err != nil {
		return err
	}

	block.enabled = true
	bm.blockList[ix][iy][iz] = block
	return nil
}

func (bm *Chunk) checkBlockPosition(x, y, z int) error {
	if x < 0 || y < 0 || z < 0 {
		return ErrBlockManagerNoSpace
	}
	if len(bm.blockList) <= x {
		return ErrBlockManagerNoSpace
	}
	if len(bm.blockList[x]) <= y {
		return ErrBlockManagerNoSpace
	}
	if len(bm.blockList[x][y]) <= z {
		return ErrBlockManagerNoSpace
	}

	return nil
}

func (c *Chunk) RenderChunk(blockModel rl.Model) {
	for typ, transforms := range c.instances {
		rl.DrawMeshInstanced(
			*blockModel.Meshes,
			blockMaterials[typ],
			transforms,
			len(transforms),
		)
	}
}

func (c *Chunk) configureBlock(x, y, z int) {
	b := c.GetBlock(x, y, z)
	if b == nil {
		return
	}

	if c.IsSurrounded(x, y, z) {
		b.enabled = false
	}

	b.blockType = BlockTypeDirt

	if c.IsSurroundedByCarved(x, y, z) {
		b.blockType = BlockTypeRock
	}

	if !c.HasBlock(x, y+1, z) {
		switch {
		case y > 38:
			b.blockType = BlockTypeSnow
		default:
			b.blockType = BlockTypeGras
		}
	}

	if y == 0 {
		b.blockType = BlockTypeGroud
	}

	if b.enabled && !b.carved {
		c.instances[b.blockType] = append(c.instances[b.blockType], rl.MatrixTranslate(
			b.position.X,
			b.position.Y,
			b.position.Z,
		))
	}
}

func (c *Chunk) Generate() {
	for x := 0; x < len(c.blockList); x++ {
		for y := 0; y < len(c.blockList[x]); y++ {
			for z := 0; z < len(c.blockList[x][y]); z++ {
				c.configureBlock(x, y, z)
			}
		}
	}
}
