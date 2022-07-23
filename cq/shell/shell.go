package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/term"
)

var (
	sigChannel = make(chan os.Signal, 1)
)

type Shell struct {
	isTTY bool
}

func NewShell() *Shell {
	return &Shell{
		isTTY: term.IsTerminal(int(os.Stdin.Fd())),
	}
}

func setup() {
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		select {
		case <-sigChannel:
			os.Exit(0)
		}
	}()
}

// func readPrivate(prompt string) (string, error) {
// 	fd := int(os.Stdin.Fd())
// 	fmt.Print(prompt)
// 	oldState, err := term.MakeRaw(fd)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to stash terminal: %w", err)
// 	}

// 	t := term.NewTerminal(os.Stdin, "")
// 	for {
// 		response, err := t.ReadPassword("")
// 		if err != nil {
// 			return "", fmt.Errorf("unable to read password: %w", err)
// 		}
// 		response = strings.TrimSpace(response)
// 		if err := term.Restore(fd, oldState); err != nil {
// 			return "", fmt.Errorf("failed to restore terminal: %w", err)
// 		}
// 		fmt.Print("\n")
// 		return response, nil
// 	}
// }

func (s *Shell) Start() error {
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
	switch input {
	//case "password":
	//readPrivate("Enter password:")
	default:
		fmt.Printf("--%s\n", input)
	}
	return nil
}
