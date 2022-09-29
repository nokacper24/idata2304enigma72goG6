package main

import (
	"fmt"
	"net"
	"strings"
	"time"
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

	hdwq, err := processRespone(returnMessage)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error processing UDP response")
		return
	}

	fmt.Println(hdwq)

	conn.Write([]byte(hdwq))

	secondResponse, err := readUDPResp(conn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error reading UDP response")
		return
	}

	fmt.Println(secondResponse)

	defer conn.Close()
}

/**
Creates a UDP connection to the server and sends the initial connection message

@return net.Conn - the connection to the server
*/
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

/**
Reads the response from the server.

@param net.Conn - the connection to the server

@return string - the response from the server
*/
func readUDPResp(conn net.Conn) (string, error) {
	// Read response from server
	buf := make([]byte, 1024)
	time.Sleep(1 * time.Second)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

/**
Processes the response from the server and returns the correct response.

@return string - the correct response to send back to the server
*/
func processRespone(response string) (string, error) {
	var responseString string
	// Process response from server
	// get the last character of the response
	punctuation := response[len(response)-1:]

	var quorsta string
	if punctuation == "." {
		quorsta = "statement"
	}
	if punctuation == "?" {
		quorsta = "question"
	}

	//get the amount of words in the response
	words := strings.Fields(response)
	wordCount := len(words)

	//get the amount of characters in the response
	charCount := len(response)

	if charCount == 1 {
		wordCount = 0
	}

	wordCountString := fmt.Sprintf("%d", wordCount)

	responseString = quorsta + " " + wordCountString

	return responseString, nil
}
