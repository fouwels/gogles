package main // import "github.com/kaelanfouwels/gogles"

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/kaelanfouwels/gogles/ioman"
	"github.com/kaelanfouwels/gogles/mfdman"

	"github.com/kaelanfouwels/gogles/fontman"
	"github.com/kaelanfouwels/gogles/textman"

	//"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kaelanfouwels/gogles/glfw/v3.3/glfw"
	"github.com/kaelanfouwels/gogles/glow/gl"

	//gl "github.com/kaelanfouwels/gogles/glow/gles"
	"github.com/kaelanfouwels/gogles/renderman"

	"flag"
)

const _width, _height = 800, 480
const _glLoopTime = (1 * time.Second) / 60       // 60 Hz
const _processLoopTime = (1 * time.Second) / 120 // 120 Hz

var flagNoGui *bool

func init() {
	//GLFW event handling must run on the main OS thread
	logf("init", "Locking to OS Thread")
	runtime.LockOSThread()

	//Commandline Flags
	logf("init", "Parsing Flags")
	flagNoGui = flag.Bool("no-gui", false, "run application in headless (no GUI) mode")
	flag.Parse()
}

func main() {
	err := start()
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

func start() error {

	logf("start", "Initializing ioman")
	ioman, err := ioman.NewIOMan()
	if err != nil {
		return err
	}
	defer ioman.Destroy()

	chioerr := make(chan error)
	logf("start", "Starting watchdog")
	go watchdog(chioerr)

	logf("start", "Starting ioman")
	go ioman.Start(chioerr)

	if !*flagNoGui {

		logf("start", "Handing over to graphics (glloop) at %v hz", 1/_glLoopTime.Seconds())
		gltick := time.NewTicker(_glLoopTime)
		defer gltick.Stop()

		err := graphics(gltick.C, ioman)
		if err != nil {
			return fmt.Errorf("glloop has exit: %w", err)
		}

	} else {
		logf("start", "Running in headless mode, handing over to CLI (cliloop)")

		cltick := time.NewTicker(1 * time.Second)
		defer cltick.Stop()

		err := cli(cltick.C, ioman)
		if err != nil {
			return fmt.Errorf("cliloop has exit: %w", err)
		}
	}

	return fmt.Errorf("glLoop exit without error, this is unexpected")
}

func watchdog(ioman <-chan error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case err := <-ioman:
			logf("watchdog", "Ioman has raised fault, exiting: %v", err)
			os.Exit(1)
		default:
		}
	}
}
func cli(ticker <-chan time.Time, ioman *ioman.IOMan) error {

	logf("cli", "Starting")
	for range ticker {
		dp := ioman.GetDataPacket()
		logf("cli", "%v", dp)
	}

	return fmt.Errorf("cli has exit unexpectedly")
}

func graphics(ticker <-chan time.Time, ioman *ioman.IOMan) error {

	logf("graphics", "Initializing GLFW")
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	logf("graphics", "Requesting Window")
	window, err := glfw.CreateWindow(_width, _height, "gogles", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	logf("graphics", "Initializing texman")
	textman, err := textman.NewTextman("./assets")
	if err != nil {
		return err
	}
	defer textman.Destroy()

	logf("graphics", "Initializing fontman")
	fontman, err := fontman.NewFontman(textman)
	if err != nil {
		return err
	}

	logf("graphics", "Initializing mdfman")
	mfdman, err := mfdman.NewMFDman(_width, _height, fontman)
	if err != nil {
		return err
	}

	logf("graphics", "Initializing renderman")
	renderman, err := renderman.NewRenderman(_width, _height, textman, fontman, mfdman, ioman)
	if err != nil {
		return err
	}
	defer renderman.Destroy()

	logf("graphics", "Starting Draw Cycle")
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
		err = fontman.RenderString(fmt.Sprintf("Healthkeeper v0.1: %v", ticks), -_width/2+20, -_height/2+20, 0.10)
		if err != nil {
			return err
		}
		ticks++
		//DEBUG

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return fmt.Errorf("graphis has exit unexpectedly")
}

func processLoop(ticker <-chan time.Time, cherr chan<- error) {

	for range ticker {

	}
}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
