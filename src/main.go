package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	var debug bool
	var silent bool
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.BoolVar(&silent, "silent", false, "silent")
	flag.Parse()

	tailLogger := tail.DiscardingLogger
	if debug {
		tailLogger = tail.DefaultLogger
	}
	tailLogger.SetPrefix("tail-debounce-exec-greps ")
	log.SetPrefix("tail-debounce-exec-greps ")

	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	if len(flag.Args()) < 4 {
		fmt.Println("usage: tail-debounce-exec-greps file debounceDuration command grep...")
		os.Exit(1)
	}

	filename := flag.Args()[0]
	duration, err := time.ParseDuration(flag.Args()[1])
	if err != nil {
		panic(err)
	}
	cmd := strings.Fields(flag.Args()[2])
	greps := flag.Args()[3:]

	log.Println(filename, duration, cmd, greps)

	t, _ := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Follow:    true,
		Logger:    tailLogger,
	})

	lastExec := time.Now()
	for line := range t.Lines {
		log.Printf(line.Text)
		for _, grep := range greps {
			if strings.Contains(line.Text, grep) {
				log.Printf("matches %s", grep)
				if time.Now().Sub(lastExec) < duration {
					log.Printf("debounce")
					continue
				}
				lastExec = time.Now()

				cmd := exec.Command(cmd[0], cmd[1:]...)
				if !silent {
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
				}

				switch err := cmd.Run(); err.(type) {
				case nil, *exec.ExitError:
					log.Printf("exec err %v", err)
				default:
					panic(err)
				}
			}
		}
	}
}
