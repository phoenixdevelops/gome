package common

import (
	"fmt"
	"gitlocal/gome"
	"gitlocal/gome/common/graphics"
	"go/build"
	"os"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

/*
	RenderComponent
*/

// A RenderComponent is a component used to render a texture of its
// entity
type RenderComponent struct {
	OBJPath      string
	ModelUpdated bool

	array   graphics.VertexArray
	texture uint32
}

func (rc *RenderComponent) Name() string { return "Render" }

/*
	RenderSystem
*/

// A RenderSystem renders the texture of its entities
type RenderSystem struct {
	gome.MultiSystem

	graphics.Shader
	cameraSystem *CameraSystem
	lightSystem  *LightSystem
}

func (*RenderSystem) RequiredComponents() []string { return []string{"Render", "Space"} }

func (rs *RenderSystem) Init(scene *gome.Scene) {
	// initialize the base system
	rs.MultiSystem.Init(scene)

	// initialize OpenGL
	gl.Init()

	// Configure global opengl settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 0) // set the clear color

	// if debug is enabled, show debug output
	if scene.WindowArgs.Debug {
		// opengl version
		fmt.Println("OpenGL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

		// error and debug outptut
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(func(
			source uint32,
			gltype uint32,
			id uint32,
			severity uint32,
			length int32,
			message string,
			userParam unsafe.Pointer) {

			// warn if it's an error
			errWarning := ""
			if gltype == gl.DEBUG_TYPE_ERROR {
				errWarning = "** ERROR **"
			}

			fmt.Printf("GL CALLBACK: %s type = 0x%x, severity = 0x%x, message = %s\n",
				errWarning, gltype, severity, message)
		}, gl.Ptr(nil))
	}

	// init shader
	// TODO change this hacky code
	shaderFile, err := os.Open(build.Default.GOPATH + "/src/github.com/lbuchli/gome/common/graphics/default.shader")
	if err != nil {
		panic("Could not find shader file")
	}
	rs.Shader.Init(shaderFile)

	// get the camera system, and if there isn't one, add a new instance to the scene.
	if scene.HasSystem("Camera") {
		rs.cameraSystem = scene.GetSystem("Camera").(*CameraSystem)
	} else {
		rs.cameraSystem = &CameraSystem{}
		scene.AddSystem(rs.cameraSystem)

		// there's probably also no camera entity, so add it as well
		cameraEntity := &CameraEntity{}
		cameraEntity.New()
		scene.AddEntity(cameraEntity)
	}

	// add a LightSystem if there isn't one
	if scene.HasSystem("Light") {
		rs.lightSystem = scene.GetSystem("Light").(*LightSystem)
	} else {
		rs.lightSystem = &LightSystem{}
		scene.AddSystem(rs.lightSystem)
	}
}

func (rs *RenderSystem) Add(id uint, components []gome.Component) {
	rs.MultiSystem.Add(id, components)

	renderComponent := components[0].(*RenderComponent)

	f, err := os.Open(renderComponent.OBJPath)
	if err != nil {
		gome.Throw(err, "Could not open OBJ file "+renderComponent.OBJPath)
	}

	reader := &graphics.OBJFileReader{}
	va, texture, err := reader.Data(f)
	if err != nil {
		gome.Throw(err, "Could not read data from object file")
	}

	renderComponent.array = va
	renderComponent.texture = texture
}

func (rs *RenderSystem) Update(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // apply clear color

	gl.UseProgram(rs.Shader.Program)

	// set the light uniforms
	//lightSources := rs.lightSystem.getLightSources()
	//rs.Shader.SetUniformBlock("u_Lights", lightSources, 12*4*4)

	// Projection View Matrix
	PVM := rs.cameraSystem.projectionViewMatrix()

	for _, components := range rs.MultiSystem.Entities {
		renderComponent := components[0].(*RenderComponent)
		spaceComponent := components[1].(*SpaceComponent)
		VAO := &renderComponent.array

		MVP := PVM.Mul4(spaceComponent.modelMatrix())
		rs.Shader.SetUniformFMat4("u_MVP", MVP)

		gl.BindTexture(gl.TEXTURE_2D, renderComponent.texture)

		VAO.Draw()
	}
}

func (*RenderSystem) Name() string { return "Render" }
