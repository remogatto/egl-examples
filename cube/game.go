package main

import (
	"git.tideland.biz/goas/loop"
	"github.com/remogatto/egl"
	"github.com/remogatto/egl-examples/cube/cubelib"
	"github.com/remogatto/egl/platform"
	"runtime"
	"time"
)

const (
	FRAMES_PER_SECOND = 30
	TEXTURE_PNG       = "texture/marmo.png"
)

type controlCh struct {
	eglState chan *platform.EGLState
	exit     chan bool
}

func newControlCh() *controlCh {
	return &controlCh{
		eglState: make(chan *platform.EGLState, 1),
		exit:     make(chan bool),
	}
}

// A render state includes informations about the 3d world and the EGL
// state (rendering surfaces, etc.)
type renderState struct {
	eglState *platform.EGLState
	world    *cubelib.World
	cube     *cubelib.Cube
	angle    float32
}

func (state *renderState) init(eglState *platform.EGLState) error {
	state.eglState = eglState

	display := eglState.Display
	surface := eglState.Surface
	context := eglState.Context
	width := eglState.SurfaceWidth
	height := eglState.SurfaceHeight

	if ok := egl.MakeCurrent(display, surface, surface, context); !ok {
		return egl.NewError(egl.GetError())
	}

	// Create and setup the 3D world
	state.world = cubelib.NewWorld(width, height)
	state.world.SetCamera(0.0, 0.0, 5.0)

	state.cube = cubelib.NewCube()

	if err := state.cube.AttachTextureFromFile(TEXTURE_PNG); err != nil {
		return err
	}

	state.world.Attach(state.cube)
	state.angle = 0.0

	return nil
}

// Run runs renderLoop. The loop renders a frame and swaps the buffer
// at each tick received.
func renderLoopFunc(controlCh *controlCh) loop.LoopFunc {
	return func(loop loop.Loop) error {

		var state renderState

		// Lock/unlock the loop to the current OS thread. This is
		// necessary because OpenGL functions should be called from
		// the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// We don't have yet a proper rendering state so the
		// ticker should be stopped as soon as it is created.
		ticker := time.NewTicker(time.Duration(int(time.Second) / int(FRAMES_PER_SECOND)))
		ticker.Stop()

		for {
			select {
			// At each tick render a frame and swap buffers.
			case <-ticker.C:
				state.angle += 0.05
				state.cube.RotateY(state.angle)
				state.world.Draw()
				egl.SwapBuffers(state.eglState.Display, state.eglState.Surface)

				// Receive an EGL state from the
				// native graphics subsystem and
				// initialize a rendering state.
			case eglState := <-controlCh.eglState:
				if err := state.init(eglState); err != nil {
					panic(err)
				}
				// Now that we have a proper rendering
				// state we can start the ticker.
				ticker = time.NewTicker(time.Duration(int(time.Second) / int(FRAMES_PER_SECOND)))

			case <-controlCh.exit:
				go loop.Stop()

			case <-loop.ShallStop():
				ticker.Stop()
				egl.DestroySurface(state.eglState.Display, state.eglState.Surface)
				egl.DestroyContext(state.eglState.Display, state.eglState.Context)
				egl.Terminate(state.eglState.Display)
				return nil
			}
		}
	}
}
