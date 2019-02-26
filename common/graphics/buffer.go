package graphics

import "github.com/go-gl/gl/v4.6-core/gl"

type ElementType uint32

const (
	VEC2 = ElementType(iota)
	VEC3
	VEC4
	FVEC2
	FVEC3
	FVEC4
)

// getByteSize returns the type size in bytes
func (et ElementType) getByteSize() (size int) {
	switch et {
	case VEC2:
		return 4 * 2
	case VEC3:
		return 4 * 3
	case VEC4:
		return 4 * 4
	case FVEC2:
		return 4 * 2
	case FVEC3:
		return 4 * 3
	case FVEC4:
		return 4 * 4
	default:
		return 0
	}
}

// getGLType returns the type as a OpenGL constant
func (et ElementType) getGLType() (gltype uint32) {
	switch et {
	case VEC2:
		return gl.INT_VEC2
	case VEC3:
		return gl.INT_VEC3
	case VEC4:
		return gl.INT_VEC4
	case FVEC2:
		return gl.FLOAT_VEC2
	case FVEC3:
		return gl.FLOAT_VEC3
	case FVEC4:
		return gl.FLOAT_VEC4
	default:
		return 0
	}
}

// A VertexLayout contains the buffer layout of a Vertex.
type VertexLayout struct {
	layout []ElementType
}

// Push pushes a type onto the VertexBufferLayout.
func (vbl *VertexLayout) Push(eType ElementType) {
	vbl.layout = append(vbl.layout, eType)
}

// stride returns the stride size of the layout.
func (vbl *VertexLayout) stride() int {
	size := 0
	for _, i := range vbl.layout {
		size += i.getByteSize()
	}

	return size
}

// A VertexArray is an array of vertices saved on the GPU memory.
type VertexArray struct {
	layout VertexLayout
	data   []interface{}
	vao    uint32
	ibo    uint32
	vbos   []uint32
}

// SetLayout sets the vertex layout, but only once.
func (va *VertexArray) SetLayout(layout VertexLayout) {
	if len(va.layout.layout) != 0 {
		return
	}

	va.layout = layout

	// generate and bind the vertex array
	gl.GenVertexArrays(1, &va.vao) // generates the vertex array (or multiple)
	gl.BindVertexArray(va.vao)     // binds the vertex array

	// make vertex array pointer attributes
	offset := 0

	// calculate vertex stride
	stride := 0
	for _, elem := range va.layout.layout {
		stride += elem.getByteSize()
	}

	for i, elem := range va.layout.layout {
		// define an array of generic vertex attribute data
		// index, size, type, normalized, stride of vertex (in bytes), pointer (offset)
		// point positions
		gl.VertexAttribPointer(uint32(i), int32(elem.getByteSize()),
			elem.getGLType(), false, int32(stride), gl.Ptr(uint32(offset)))
		gl.EnableVertexAttribArray(uint32(i))
		offset += elem.getByteSize()
	}

	// make as many vbos as there are vertex array pointer attributes
	va.vbos = make([]uint32, len(layout.layout))
}

// SetData sets the buffer data at a specific index to be equal to the slice of data.
func (va *VertexArray) SetData(index int, data interface{}) (err error) {
	// Vertex Buffer Object
	var VBO uint32
	gl.GenBuffers(1, &VBO)              // generates the buffer (or multiple)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO) // tells OpenGL what kind of buffer this is

	// BufferData assigns data to the buffer.
	// there can only be one ARRAY_BUFFER, so OpenGL knows which buffer we mean if we
	// tell it what type of buffer it is.
	//			  type			   size (in bytes)   pointer to data	usage
	gl.BufferData(gl.ARRAY_BUFFER, 0, gl.Ptr(0), gl.STATIC_DRAW)

	va.vbos[index] = VBO

	return
}

// SetIndexData sets the index buffer object of the array.
func (va *VertexArray) SetIndexData(data []uint32) {
	// Index Buffer Object
	gl.GenBuffers(1, &va.ibo)                      // generates the buffer (or multiple)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, va.ibo) // tells OpenGL what kind of buffer this is

	// BufferData assigns data to the buffer.
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
}

// Bind binds the VertexArray, allowing it to be drawn.
func (va *VertexArray) Bind() {
	gl.BindVertexArray(va.vao)
}