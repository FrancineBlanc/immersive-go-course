package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	exitCode := client()
	if exitCode == 1 {
		os.Exit(1)
	}
	// TO CHECK - What if there was more code in main, and an err has been thrown
}

func getRequest() (string, string, int, error) {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		return "", "", 0, fmt.Errorf("the weather cannot be retrieved right now")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, fmt.Errorf("unable to read weather data at present")
	}
	return string(body), resp.Header.Get("Retry-After"), resp.StatusCode, nil
}

func retryAfterAWhile() {
	fmt.Println("Experiencing high volumes of traffic...retrying again...")
	time.Sleep(time.Duration(4) * time.Second)
	client()
}

func retryInSeconds(secs int) {
	fmt.Println("Retrying in", secs, "second(s)...")
	time.Sleep(time.Duration(secs) * time.Second)
	client()
}

func calculateSecs(retryHeader string) int {
	secs, err := strconv.Atoi(retryHeader)
	if err != nil {
		t, _ := http.ParseTime(retryHeader)
		unixTime := t.UTC().Unix()
		currentTime := time.Now().UTC().Unix()
		secs = int(unixTime - currentTime)
	}
	return secs
}

func client() int {
	body, retryAfterHeader, statusCode, err := getRequest()
	if err != nil {
		os.Stderr, _ = os.Create("stderr-log.txt")
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "Connection dropped: %v\n", err)
		return 1
	}
	if statusCode == 200 {
		fmt.Println(string(body))
	} else if statusCode == 429 {
		if retryAfterHeader == "a while" {
			retryAfterAWhile()
		} else {
			secs := calculateSecs(retryAfterHeader)
			if secs >= 1 && secs < 5 {
				retryInSeconds(secs)
			} else {
				fmt.Println("Weather cannot be retrieved :(")
			}
		}
	}
	return 0
}
