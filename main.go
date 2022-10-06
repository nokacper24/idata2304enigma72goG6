package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// import env variables package
	"github.com/joho/godotenv"

	//import a logging package
	"socketProgrammingUDP/logger"

	//import library to make a web gui for the server
	"github.com/gin-gonic/gin"
)


var log = logger.NewLogger()

func main() {
	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
		panic(err)
	}

	// get the url and port from the env variables
	url := os.Getenv("URL")
	port := os.Getenv("PORT")
	message := os.Getenv("MESSAGE")


	// get the number of tasks to perform from the env variables
	tasks, err := strconv.Atoi(os.Getenv("TASKS"))
	if err != nil {
		log.Error("Error converting TASKS to int")
		panic(err)
	}

	webURL := os.Getenv("WEB_URL")
	webPort := os.Getenv("WEB_PORT")

	// create a new router
	router := gin.Default()

	// create a new group of routes
	frontend := router.Group("/")

	// create a route to serve the index.html file
	frontend.StaticFile("/", "./assets/index.html")

	// create a route to serve the javascript file
	frontend.StaticFile("/script.js", "./assets/script.js")

	// create a route to serve the css file
	frontend.StaticFile("/style.css", "./assets/style.css")


	
	frontend.GET("/tasks", func(c *gin.Context) {
		msg,tsk := performTask(url, port, message)
		if tsk {
		c.JSON(200, gin.H{
			"tasks": msg,
		})
		}else{
			c.JSON(500, gin.H{
				"tasks": msg,
			})
		}
	})
	
	router.Run(webURL + ":" + webPort)

	//print env variables
	fmt.Println("URL: " + url)
	fmt.Println("PORT: " + port)
	fmt.Println("MESSAGE: " + message)
	fmt.Println("TASKS: " + strconv.Itoa(tasks))


	numberOfWorkers := tasks

	var wg sync.WaitGroup

	var successfullConnections int
	var failedConnections int

	// log time taken to perform tasks
	start := time.Now()

	for i := 0; i <= numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			returnString, shouldReturn := performTask(url, port, message)
			if !shouldReturn {
				// if the task fails, print the error with red bg and return
				fmt.Println("\033[41m" + "Error performing task" + "\033[0m")
				failedConnections++
				wg.Done()

			}
			if shouldReturn {
				successfullConnections++
				log.Info(returnString)
				wg.Done()
			}

		}()

	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Time taken to perform tasks: " + elapsed.String())
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
func performTask(url string, port string, message string) ([]string,bool) {
	// store the return messages so they can be printed at the end
	var returnMessages []string
	conn, err := sendUDPReq(url, port, message)
	if err != nil {
		log.Error("Error sending UDP request")
		log.Error(err)
		return returnMessages,false
	}
	initialQuestion, err := readUDPResp(conn)
	if err != nil {
		log.Error("Error reading UDP response")
		log.Error(err)
		return returnMessages,false
	}

	// print a string of ---- to separate the initial question from the response
	returnMessages = append(returnMessages, initialQuestion)

	

	questionAnswer, err := processRespone(initialQuestion)
	if err != nil {
		log.Error("Error processing response")
		log.Error(err)
		return returnMessages,false
	}

	returnMessages = append(returnMessages, questionAnswer)

	conn.Write([]byte(questionAnswer))

	correctResponseAnswer, err := readUDPResp(conn)
	if err != nil {
		log.Error("Error reading UDP response")
		log.Error(err)
		return returnMessages,false
	}

	if correctResponseAnswer == "ok" {
		//print the response from the server 
		returnMessages = append(returnMessages, correctResponseAnswer)
	} else {
		//print the response from the server 
		returnMessages = append(returnMessages, correctResponseAnswer)
	}

	// prepare the return messages to be printed in a chunk so they are not incorrectly ordered


	defer conn.Close()
	return returnMessages,true
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
