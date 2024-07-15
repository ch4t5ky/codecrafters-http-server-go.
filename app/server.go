package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running test
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading: %#v\n", err)
			return
		}
		response := HandleConnection(string(buf))
		conn.Write(response)
		conn.Close()
	}
}

func HandleConnection(request string) []byte {
	path := strings.Split(request, " ")[1]
	response := []byte("")
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if path == "/user-agent" {
		msg := ""
		packet := strings.Split(request, "\r\n")
		for i := 0; i < len(packet); i++ {
			fmt.Println(packet[i])
			dict := strings.Split(packet[i], ":")
			if len(dict) != 2 {
				continue
			}
			header, value := dict[0], dict[1]
			header = strings.ToLower(header)
			if header == "user-agent" {
				msg = strings.ReplaceAll(value, " ", "")
				break
			}
		}

		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg))
	} else if strings.Split(path, "/")[1] == "echo" {
		msg := strings.Split(path, "/")[2]
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg))
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	return response
}
