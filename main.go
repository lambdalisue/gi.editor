package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/comail/colog"
)

var (
	appVersion = "dev"
)

func main() {
	var (
		version = flag.Bool("version", false, "show version")
		addr    = flag.String("addr", "", "TCP address to listen")
	)
	flag.Parse()
	colog.Register()

	if *version {
		fmt.Println(appVersion)
		os.Exit(0)
	} else if *addr == "" {
		log.Fatalf("error: -addr must be specified\n")
	} else if flag.NArg() < 1 {
		log.Fatalf("error: no file has specified\n")
	}

	exitCode, err := run(*addr, flag.Args()[0])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	os.Exit(exitCode)
}

func run(addr, file string) (int, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return 1, fmt.Errorf("faield to dial %s: %w", addr, err)
	}
	defer conn.Close()

	// Notify the filename to the server
	w := bufio.NewWriter(conn)
	if _, err := w.WriteString(file + "\n"); err != nil {
		return 1, fmt.Errorf("failed to write %s: %w", conn, err)
	}
	if err := w.Flush(); err != nil {
		return 1, fmt.Errorf("failed to flush %s: %w", conn, err)
	}

	// Read FIRST line from the server as exitCode and exit
	r := bufio.NewReader(conn)
	recv, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF || err == io.ErrClosedPipe {
			return 1, fmt.Errorf("exitCode must be given prior to EOF")
		}
		return 1, fmt.Errorf("failed to read: %w", err)
	}
	exitCode, err := strconv.Atoi(recv)
	if err != nil {
		return 1, fmt.Errorf("failed to parse %s: %w", recv, err)
	}
	return exitCode, nil
}
