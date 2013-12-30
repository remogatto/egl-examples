package main

import (
	"flag"
	"fmt"
	"git.tideland.biz/goas/loop"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
)

func main() {
	size := flag.String("size", "320x480", "set the size of the window")

	flag.Parse()

	dims := strings.Split(strings.ToLower(*size), "x")
	width, err := strconv.Atoi(dims[0])
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(dims[1])
	if err != nil {
		panic(err)
	}

	controlCh := newControlCh()
	controlCh.eglState <- initEGL(controlCh, width, height)

	// Start the rendering loop
	loop.GoRecoverable(
		renderLoopFunc(controlCh),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				log.Println(r.Reason)
				log.Println(string(debug.Stack()))
			}
			return rs, fmt.Errorf("Unrecoverable loop error\n")
		},
	).Wait()
}
