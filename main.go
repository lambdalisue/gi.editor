package main

import (
	"fmt"
	"net"
	"os"
	"flag"
	"strconv"
)

var (
	appVersion = "dev"
)

const (
	exitFatal = 1
	exitFatalArgs = 3
	exitFatalListen
	exitFatalAccept
	exitFatalRead
	exitFatalParse
)

func run(addr string) (int, error) {
	args := flag.Args()
	if len(args) < 1 {
		return exitFatalArgs, fmt.Errorf("No file has specified.")
	}
	fmt.Println("")
	fmt.Println("file:" + args[0])
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return exitFatalListen, err
	}
	fmt.Println("addr:" + l.Addr().String())
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					continue
				}
			}
			return exitFatalAccept, err
		}
		// Accept the first request only
		defer conn.Close()
		var b [16]byte
		n, err := conn.Read(b[:])
		if err != nil {
			return exitFatalRead, err
		}
		exitCode, err := strconv.Atoi(string(b[:n]))
		if err != nil {
			return exitFatalParse, err
		}
		return exitCode, nil
	}
}

func main() {
	var (
		version = flag.Bool("version", false, "show version")
		addr = flag.String("addr", "127.0.0.1:0", "TCP address to listen")
	)
	flag.Parse()
	if *version {
		fmt.Println(appVersion)
		os.Exit(exitFatal)
	}
	exitCode, err := run(*addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}
