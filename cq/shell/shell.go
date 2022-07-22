package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
)

var (
	sigChannel = make(chan os.Signal, 1)
)

func setup() {
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		select {
		case <-sigChannel:
			os.Exit(0)
		}
	}()
}

func Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("cq > ")

		line, err := reader.ReadString('\n')
		switch {
		case errors.Is(err, io.EOF):
			os.Exit(0)
		case err != nil:
			fmt.Fprintln(os.Stderr, err)
		}

		switch {
		case len(line) >= 2 && line[len(line)-2] == '\r':
			line = line[:len(line)-2]
		default:
			line = line[:len(line)-1]
		}

		if err = execInput(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	fmt.Println(input)
	return nil
}
