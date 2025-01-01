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

func main() {
	append := flag.Bool("a", false, "append to the given file")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	outName := flag.Arg(0)
	tempFile, err := os.CreateTemp("", "sponge.*.tmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
		os.Exit(1)
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

	buf := bufio.NewReader(os.Stdin)
	if _, err := io.Copy(tempFile, buf); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read from stdin: %v\n", err)
		os.Exit(1)
	}

	if err := tempFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close temp file: %v\n", err)
		os.Exit(1)
	}

	if *append {
		outFile, err := os.OpenFile(outName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()

		tempFile, err := os.Open(tempFile.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to reopen temp file: %v\n", err)
			os.Exit(1)
		}
		defer tempFile.Close()

		if _, err := io.Copy(outFile, tempFile); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := os.Rename(tempFile.Name(), outName); err != nil {
			outFile, err := os.OpenFile(outName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
				os.Exit(1)
			}
			defer outFile.Close()

			tempFile, err := os.Open(tempFile.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to reopen temp file: %v\n", err)
				os.Exit(1)
			}
			defer tempFile.Close()

			if _, err := io.Copy(outFile, tempFile); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write to output file: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
