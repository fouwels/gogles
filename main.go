package main // import "github.com/kaelanfouwels/gogles"

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/kaelanfouwels/gogles/mfdman"

	"github.com/kaelanfouwels/gogles/fontman"
	"github.com/kaelanfouwels/gogles/textman"

	"github.com/go-gl/glfw/v3.3/glfw"
	gl "github.com/kaelanfouwels/gogles/glow/gl"
	"github.com/kaelanfouwels/gogles/renderman"
	//gl "github.com/kaelanfouwels/gogles/glow/gles"
)

const width, height = 800, 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	err := start()
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

func start() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, "gogles", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	textman, err := textman.NewTextman("./assets")
	if err != nil {
		return err
	}
	defer textman.Destroy()

	fontman, err := fontman.NewFontman(textman)
	if err != nil {
		return err
	}

	mfdman, err := mfdman.NewMFDman(width, height, fontman)
	if err != nil {
		return err
	}

	renderman, err := renderman.NewRenderman(width, height, textman, fontman, mfdman)
	if err != nil {
		return err
	}
	defer renderman.Destroy()

	go simulatedMFD(mfdman)

	ticks := 0
	for !window.ShouldClose() {

		err := renderman.Draw()
		if err != nil {
			return fmt.Errorf("Draw cycle failed: %w", err)
		}
		err = fontman.RenderString(fmt.Sprintf("Healthkeeper v0.1: %v", ticks), -width/2+20, -height/2+20, 0.10)
		if err != nil {
			return err
		}
		window.SwapBuffers()
		glfw.PollEvents()
		ticks++
		time.Sleep(10 * time.Millisecond)

	}

	return nil
}

func simulatedMFD(mman *mfdman.MFDman) {

	mman.SetText(mfdman.L4, "-", "FLOW")
	mman.SetText(mfdman.R4, "+", "FLOW")
	mman.SetText(mfdman.L3, "-", "O2")
	mman.SetText(mfdman.R3, "+", "O2")

	mman.SetText(mfdman.R1, "MENU", "MENU")

	for true {

		for i := mfdman.L1; i < mfdman.MFDCount; i++ {
			time.Sleep(500 * time.Millisecond)
			mman.SetSelected(i, !mman.GetSelected(i))
			time.Sleep(500 * time.Millisecond)
			mman.SetSelected(i, !mman.GetSelected(i))
		}
	}
}
