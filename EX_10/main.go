package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// Парсинг аргументов командной строки
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		return
	}

	host := flag.Args()[0]
	port := flag.Args()[1]

	address := net.JoinHostPort(host, port)

	// Установка соединения
	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	// Канал для завершения программы
	done := make(chan struct{})

	// Чтение из сокета и вывод в STDOUT
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Println("Connection closed by server")
		done <- struct{}{}
	}()

	// Чтение из STDIN и запись в сокет
	go func() {
		io.Copy(conn, os.Stdin)
		conn.Close()
		done <- struct{}{}
	}()

	// Ожидание завершения
	<-done
}
