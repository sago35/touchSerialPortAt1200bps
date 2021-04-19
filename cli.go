package main

import (
	"fmt"
	"io"
	"runtime"
	"time"

	"go.bug.st/serial"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName        = "touchSerialPortAt1200bps"
	appDescription = ""
)

type cli struct {
	outStream io.Writer
	errStream io.Writer
}

var (
	app  = kingpin.New(appName, appDescription)
	port = app.Arg("port", "flash port").Required().String()
)

// Run ...
func (c *cli) Run(args []string) error {
	app.UsageWriter(c.errStream)

	if VERSION != "" {
		app.Version(fmt.Sprintf("%s version %s build %s", appName, VERSION, BUILDDATE))
	} else {
		app.Version(fmt.Sprintf("%s version - build -", appName))
	}
	app.HelpFlag.Short('h')

	k, err := app.Parse(args[1:])
	if err != nil {
		return err
	}

	switch k {
	default:
		err := touchSerialPortAt1200bps(*port)
		if err != nil {
			return err
		}
	}

	return nil
}

func touchSerialPortAt1200bps(port string) (err error) {
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		// Open port
		p, e := serial.Open(port, &serial.Mode{BaudRate: 1200})
		if e != nil {
			if runtime.GOOS == `windows` {
				se, ok := e.(*serial.PortError)
				if ok && se.Code() == serial.InvalidSerialPort {
					// InvalidSerialPort error occurs when transitioning to boot
					return nil
				}
			}
			time.Sleep(1 * time.Second)
			err = e
			continue
		}
		defer p.Close()

		p.SetDTR(false)
		return nil
	}
	return fmt.Errorf("opening port: %s", err)
}
