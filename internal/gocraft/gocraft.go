package gocraft

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/urfave/cli/v2"
)

type engine struct {
	// window settings
	screenWidth  int32
	screenheight int32
	title        string
	loglevel     rl.TraceLogLevel
	maxFPS       int
	camera       rl.Camera3D
	cunkMan      *ChunkManager

	// cam vars (move to own object)
	cameraDirection rl.Vector3
	cameraRight     rl.Vector3
	cameraFront     rl.Vector3
}

func newEngine(ctx *cli.Context) *engine {
	return &engine{
		screenWidth:  int32(ctx.Int("width")),
		screenheight: int32(ctx.Int("height")),
		title:        ctx.String("title"),
		loglevel:     loglevelFromString(ctx.String("loglevel")),
		maxFPS:       ctx.Int("fps"),
		camera: rl.NewCamera3D(
			rl.NewVector3(0, 20, 0),
			rl.NewVector3(0, 100, 100),
			rl.NewVector3(0, 1, 0),
			60,
			rl.CameraPerspective,
		),
		cameraFront: rl.NewVector3(0, 0, -1),
		cunkMan:     NewChunkManager(64),
	}
}

func RunEngine(ctx *cli.Context) error {
	es := newEngine(ctx)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(es.screenWidth, es.screenheight, es.title)
	rl.SetTargetFPS(int32(es.maxFPS))

	mainLoop(es)

	rl.CloseWindow()
	return nil
}

func mainLoop(state *engine) {
	cube := rl.LoadModel("res/models/brick.obj")
	// Load basic lighting shader
	shader := rl.LoadShader("res/shaders/lighting.vs", "res/shaders/lighting.fs")
	// Get some required shader locations
	shader.UpdateLocation(rl.LocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.LocVectorView, rl.GetShaderLocation(shader, "viewPos"))
	shader.UpdateLocation(rl.LocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))

	// ambient light level
	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	rl.SetShaderValue(shader, ambientLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
	NewLight(LightTypeDirectional, rl.NewVector3(50.0, 50.0, 0.0), rl.Vector3Zero(), rl.Beige, shader)

	materials := make([]rl.Material, 3)

	materials[0] = rl.LoadMaterialDefault()
	materials[0].Shader = shader
	materials[0].Maps.Texture = rl.LoadTexture("res/textures/dirt.png")
	materials[0].GetMap(rl.MapDiffuse).Color = rl.White

	materials[1] = rl.LoadMaterialDefault()
	materials[1].Shader = shader
	materials[1].Maps.Texture = rl.LoadTexture("res/textures/gras.png")
	materials[1].GetMap(rl.MapDiffuse).Color = rl.White

	materials[2] = rl.LoadMaterialDefault()
	materials[2].Shader = shader
	materials[2].Maps.Texture = rl.LoadTexture("res/textures/snow.png")
	materials[2].GetMap(rl.MapDiffuse).Color = rl.White

	firstChunk := state.cunkMan.GetChunk(rl.Vector2Zero(), rl.White)
	for !rl.WindowShouldClose() {
		// Update the light shader with the camera view position
		rl.SetShaderValue(shader, shader.GetLocation(rl.LocVectorView),
			[]float32{state.camera.Position.X, state.camera.Position.Y, state.camera.Position.Z}, rl.ShaderUniformVec3)

		rl.BeginDrawing()

		rl.ClearBackground(rl.SkyBlue)
		processInput(state)
		updateCamera(state)
		rl.BeginMode3D(state.camera)
		{
			chunks := state.cunkMan.GetChunks(state.camera.Position, firstChunk)
			for i, chunk := range chunks {
				if i == 0 {
					firstChunk = chunk
				}
				chunk.BuildMesh(cube, materials)
			}
			rl.DrawGrid(128, 128)
		}
		rl.EndMode3D()
		state.cunkMan.DebugChunks(state.camera.Position, firstChunk)
		rl.DrawFPS(5, 5)
		rl.EndDrawing()
	}
}

var (
	moveSpeed           = 13
	sensitivity float32 = 0.3
	yaw         float32
	pitch       float32
)

// TODO: make own object
func updateCamera(s *engine) {
	s.cameraDirection = rl.Vector3Subtract(s.camera.Position, s.camera.Target)
	s.cameraDirection = rl.Vector3Normalize(s.cameraDirection)
	s.cameraRight = rl.Vector3CrossProduct(s.camera.Up, s.cameraDirection)
	s.cameraRight = rl.Vector3Normalize(s.cameraRight)

	s.camera.Target = rl.Vector3Add(s.camera.Position, s.cameraFront)
}

func processInput(s *engine) {
	if rl.IsKeyDown(rl.KeyW) {
		s.camera.Position = rl.Vector3Add(
			s.camera.Position,
			rl.Vector3Scale(s.cameraFront, float32(moveSpeed)*rl.GetFrameTime()),
		)
	}
	if rl.IsKeyDown(rl.KeyS) {
		s.camera.Position = rl.Vector3Subtract(
			s.camera.Position,
			rl.Vector3Scale(s.cameraFront, float32(moveSpeed)*rl.GetFrameTime()),
		)
	}
	if rl.IsKeyDown(rl.KeyD) {
		s.camera.Position = rl.Vector3Add(
			s.camera.Position,
			rl.Vector3Scale(
				rl.Vector3Normalize(rl.Vector3CrossProduct(s.cameraFront, s.camera.Up)),
				float32(moveSpeed)*rl.GetFrameTime(),
			),
		)
	}
	if rl.IsKeyDown(rl.KeyA) {
		s.camera.Position = rl.Vector3Subtract(
			s.camera.Position,
			rl.Vector3Scale(
				rl.Vector3Normalize(rl.Vector3CrossProduct(s.cameraFront, s.camera.Up)),
				float32(moveSpeed)*rl.GetFrameTime(),
			),
		)
	}
	if rl.IsWindowFocused() {
		mousepos := rl.GetMouseDelta()

		yaw += mousepos.X * sensitivity
		pitch += -mousepos.Y * sensitivity

		var direction rl.Vector3

		direction.X = cos32(radians(yaw)) * cos32(radians(pitch))
		direction.Y = sin32(radians(pitch))
		direction.Z = sin32(radians(yaw)) * cos32(radians(pitch))
		s.cameraFront = rl.Vector3Normalize(direction)

		rl.SetMousePosition(int(s.screenWidth)/2, int(s.screenheight)/2)
		rl.HideCursor()
	} else {
		rl.ShowCursor()
	}
}
