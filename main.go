package main // import "github.com/kaelanfouwels/gogles"

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kaelanfouwels/gogles/ioman"
	"github.com/kaelanfouwels/gogles/server"

	"flag"
)

const _width, _height = 800, 480
const _glLoopTime = (1 * time.Second) / 60 // 60 Hz
const _cliLoopTime = (1 * time.Second) / 1 // 1 Hz

var flagNoIo *bool

func init() {

	//Commandline Flags
	logf("init", "Parsing Flags")
	flagNoIo = flag.Bool("no-io", false, "run the application without IO")
	flag.Parse()
}

func main() {
	err := start()
	if err != nil {
		logf("main", "%v", err)
		os.Exit(1)
	}
}

func start() error {

	chioerr := make(chan error)
	chserver := make(chan error)

	ioman1 := &ioman.IOMan{}
	if !*flagNoIo {
		logf("start", "Initializing ioman")
		var err error
		ioman1, err = ioman.NewIOMan()
		if err != nil {
			return err
		}
		defer ioman1.Destroy()

		logf("start", "Starting ioman")
		go ioman1.Start(chioerr)
	}

	logf("start", "Starting http data server")
	go server.Serve(chserver, ioman1)

	logf("start", "All services started, dropping into watchdog")
	watchdog(chioerr, chserver)

	return fmt.Errorf("start exit without error, this is unexpected")
}

func watchdog(ioman <-chan error, http <-chan error) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case err := <-ioman:
			logf("watchdog", "Ioman has raised fault, exiting: %v", err)
			os.Exit(1)
		case err := <-http:
			logf("watchdog", "Server has raised fault, exiting: %v", err)
			os.Exit(1)
		default:
		}
	}
}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
