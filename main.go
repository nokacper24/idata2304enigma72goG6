package main

import (
	"fmt"
	"net"

)

func main() {
	
	url := "129.241.152.12"
	port := "1234"
	message := "task"

	conn, err := sendUDPReq(url, port, message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error sending UDP request")
		return
	}
	returnMessage, err := readUDPResp(conn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error reading UDP response")
		return
	}
	fmt.Println(returnMessage)

	defer conn.Close()
}

func sendUDPReq(url string, port string, message string) (net.Conn, error) {
	// Send UDP request to server
	conn, err := net.Dial("udp", url+":"+port)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	return conn, nil
}


func readUDPResp(conn net.Conn) (string, error) {
	// Read response from server
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}