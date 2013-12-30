// +build raspberry

package cubelib

import (
	"github.com/remogatto/egl"
	"github.com/remogatto/egl/platform/raspberry"
	"fmt"
)

const (
	INITIAL_WINDOW_WIDTH  = 1920
	INITIAL_WINDOW_HEIGHT = 1080
)

func initEGL(controlCh *controlCh, width, height int) *platform.EGLState {
	egl.BCMHostInit()
	return raspberry.Initialize(
		raspberry.DefaultConfigAttributes,
		raspberry.DefaultContextAttributes,
	)
	go func() {
		fmt.Scanln()
		controlCh.exit <- true
	}
}
