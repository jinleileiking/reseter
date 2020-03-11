package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
)

var addr = flag.String("addr", "127.0.0.1:8777", "listen addr")

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	tcpLn, ok := l.(*net.TCPListener)

	if !ok {
		panic(err)
	}

	f, err := tcpLn.File()
	if err != nil {
		panic(err)
	}

	linger := syscall.Linger{
		Onoff:  1,
		Linger: 0,
	}

	if err := syscall.SetsockoptLinger(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_LINGER, &linger); err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

		result := strconv.Itoa(666) + "\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}
