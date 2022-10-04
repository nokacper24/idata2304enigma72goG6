package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {

	url := "129.241.152.12"
	port := "1234"
	message := "task"

	numberOfWorkers := 1000

	var wg sync.WaitGroup

	var successfullConnections int
	var failedConnections int

	for i := 0; i <= numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			shouldReturn := performTask(url, port, message)
			if !shouldReturn {
				// if the task fails, print the error with red bg and return
				fmt.Println("\033[41m" + "Error performing task" + "\033[0m")
				failedConnections++
				wg.Done()

			}
			if shouldReturn {
				successfullConnections++
				wg.Done()
			}

		}()

	}
	wg.Wait()
	fmt.Println("Number of successful connections: ", successfullConnections)
	fmt.Println("Number of failed connections: ", failedConnections)

}

// performs the task of sending a UDP request to the server, reading the response and sending the correct response back to the server.
//
// @param string - the url of the server
//
// @param string - the port of the server
//
// @param string - the message to send to the server
//
// @return bool - true if the task was performed successfully, false if not
func performTask(url string, port string, message string) bool {
	// store the return messages so they can be printed at the end
	var returnMessages []string
	conn, err := sendUDPReq(url, port, message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error sending UDP request")
		return false
	}
	initialQuestion, err := readUDPResp(conn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error reading UDP response")
		return false
	}

	// print a string of ---- to separate the initial question from the response
	returnMessages = append(returnMessages, "--------------------------------------------------\n")
	returnMessages = append(returnMessages, initialQuestion+"\n")

	questionAnswer, err := processRespone(initialQuestion)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error processing UDP response")
		return false
	}

	returnMessages = append(returnMessages, questionAnswer+"\n")

	conn.Write([]byte(questionAnswer))

	correctResponseAnswer, err := readUDPResp(conn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error reading UDP response")
		return false
	}

	if correctResponseAnswer == "ok" {
		//print the response from the server with green background
		returnMessages = append(returnMessages, "\033[42m"+correctResponseAnswer+"\033[0m"+"\n")
	} else {
		//print the response from the server with red background
		returnMessages = append(returnMessages, "\033[41m"+correctResponseAnswer+"\033[0m"+"\n")
	}

	// prepare the return messages to be printed in a chunk so they are not incorrectly ordered

	largeString := ""
	for _, message := range returnMessages {
		largeString += message
	}
	fmt.Println(largeString)

	defer conn.Close()
	return true
}

// Creates a UDP connection to the server and sends the initial connection message
//
// @return net.Conn - the connection to the server
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

// Reads the response from the server.
//
// @param net.Conn - the connection to the server
//
// @return string - the response from the server
func readUDPResp(conn net.Conn) (string, error) {
	// Read response from server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	// add a timeout to the connection
	return string(buf[:n]), nil
}

// Processes the response from the server and returns the correct response.
// Returned response is 'question '/'statement ' and number of words in the given sentence.
//
// @param string - the sentence in response from the server
//
// @return string - the correct response to send back to the server
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
