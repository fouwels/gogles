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

	//"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kaelanfouwels/gogles/glfw/v3.3/glfw"
	"github.com/kaelanfouwels/gogles/glow/gl"

	//gl "github.com/kaelanfouwels/gogles/glow/gles"
	"github.com/kaelanfouwels/gogles/renderman"
)

const width, height = 800, 480
const glLoopTime = (1 * time.Second) / 60       // 60 Hz
const processLoopTime = (1 * time.Second) / 120 // 120 Hz

func init() {
	// GLFW event handling must run on the main OS thread
	logf("init", "Locking to OS Thread")
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

	logf("Start", "Starting process routine at %v hz", 1/processLoopTime.Seconds())
	processLoopTicker := time.NewTicker(processLoopTime)
	defer processLoopTicker.Stop()

	processLoopError := make(chan error)
	go processLoop(processLoopTicker.C, processLoopError)

	logf("Start", "Starting goroutine watchdog")
	go watchdog(processLoopError)

	logf("Start", "Handing over to graphics routine at %v hz", 1/glLoopTime.Seconds())
	glLoopTicker := time.NewTicker(glLoopTime)
	defer glLoopTicker.Stop()

	err := glLoop(glLoopTicker.C)
	if err != nil {
		return fmt.Errorf("glLoop has exit: %w", err)
	}
	return fmt.Errorf("glLoop exit without error, this is unexpected")
}

func watchdog(process <-chan error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case err := <-process:
			logf("watchdog", "processLoop has raised error, exiting:\n > %v", err)
			os.Exit(1)
		default:
		}
	}
}

func glLoop(ticker <-chan time.Time) error {

	logf("glloop", "Initializing GLFW")
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	logf("glloop", "Requesting Window")
	window, err := glfw.CreateWindow(width, height, "gogles", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	logf("glloop", "Initializing texman")
	textman, err := textman.NewTextman("./assets")
	if err != nil {
		return err
	}
	defer textman.Destroy()

	logf("glloop", "Initializing fontman")
	fontman, err := fontman.NewFontman(textman)
	if err != nil {
		return err
	}

	logf("glloop", "Initializing mdfman")
	mfdman, err := mfdman.NewMFDman(width, height, fontman)
	if err != nil {
		return err
	}

	logf("glloop", "Initializing renderman")
	renderman, err := renderman.NewRenderman(width, height, textman, fontman, mfdman)
	if err != nil {
		return err
	}
	defer renderman.Destroy()

	logf("glloop", "Starting Draw Cycle")
	ticks := 0
	for range ticker {

		if window.ShouldClose() {
			return fmt.Errorf("Window has been closed")
		}

		err := renderman.Draw()
		if err != nil {
			return fmt.Errorf("Draw cycle failed: %w", err)
		}

		//DEBUG
		err = fontman.RenderString(fmt.Sprintf("Healthkeeper v0.1: %v", ticks), -width/2+20, -height/2+20, 0.10)
		if err != nil {
			return err
		}
		ticks++
		//DEBUG

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}

func processLoop(ticker <-chan time.Time, cherr chan<- error) {
	for range ticker {

	}
}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
