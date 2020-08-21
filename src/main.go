package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage: tail-debounce-exec-greps file debounceDuration command grep...")
		os.Exit(1)
	}

	filename := os.Args[1]
	duration, err := time.ParseDuration(os.Args[2])
	if err != nil {
		panic(err)
	}
	cmd := strings.Fields(os.Args[3])
	greps := os.Args[4:]

	t, _ := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Follow:    true,
		Logger:    tail.DiscardingLogger,
	})

	lastExec := time.Now()
	for line := range t.Lines {
		for _, grep := range greps {
			if strings.Contains(line.Text, grep) {
				if time.Now().Sub(lastExec) < duration {
					continue
				}
				lastExec = time.Now()

				cmd := exec.Command(cmd[0], cmd[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				switch err := cmd.Run(); err.(type) {
				case nil, *exec.ExitError:
				default:
					panic(err)
				}
			}
		}
	}
}
