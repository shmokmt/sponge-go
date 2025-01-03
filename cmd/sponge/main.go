package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func usage() {
	fmt.Println("sponge [-a] <file>: soak up all input from stdin and write it to <file>")
	os.Exit(0)
}

type Sponge struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func NewSponge(Stdout io.Writer, Stderr io.Writer, Stdin io.Reader) *Sponge {
	return &Sponge{
		Stdout: Stdout,
		Stderr: Stderr,
		Stdin:  Stdin,
	}
}

func (s *Sponge) Run(outFileName string, append bool) error {
	tempFile, err := os.CreateTemp("", "sponge.*.tmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a temp file: %v\n", err)
		return err
	}
	defer os.Remove(tempFile.Name())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		tempFile.Close()
		os.Remove(tempFile.Name())
		os.Exit(1)
	}()

	buf := bufio.NewReader(s.Stdin)
	if _, err := io.Copy(tempFile, buf); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from stdin: %v\n", err)
		return err
	}

	if err := tempFile.Close(); err != nil {
		fmt.Fprintf(s.Stderr, "failed to close a temp file: %v\n", err)
		return err
	}

	if append {
		outFile, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
			return err
		}
		defer outFile.Close()

		tempFile, err := os.Open(tempFile.Name())
		if err != nil {
			fmt.Fprintf(s.Stderr, "failed to reopen temp file: %v\n", err)
			return err
		}
		defer tempFile.Close()

		if _, err := io.Copy(outFile, tempFile); err != nil {
			fmt.Fprintf(s.Stderr, "failed to write to output file: %v\n", err)
			return err
		}
	} else {
		if err := os.Rename(tempFile.Name(), outFileName); err != nil {
			outFile, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Fprintf(s.Stderr, "failed to open an output file: %v\n", err)
				return err
			}
			defer outFile.Close()

			tempFile, err := os.Open(tempFile.Name())
			if err != nil {
				fmt.Fprintf(s.Stderr, "failed to reopen a temp file: %v\n", err)
				return err
			}
			defer tempFile.Close()

			if _, err := io.Copy(outFile, tempFile); err != nil {
				fmt.Fprintf(s.Stderr, "failed to write to an output file: %v\n", err)
				return err
			}
		}
	}
	return nil
}

func main() {
	append := flag.Bool("a", false, "append to the given file")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	outName := flag.Arg(0)
	sponge := NewSponge(os.Stdout, os.Stderr, os.Stdin)
	if err := sponge.Run(outName, *append); err != nil {
		os.Exit(1)
	}
}
