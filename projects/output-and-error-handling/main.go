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
	body, _, statusCode, err := getRequest()
	if err != nil {
		os.Stderr, _ = os.Create("stderr-log.txt")
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "Connection dropped: %v\n", err)
		os.Exit(1)
	} else if statusCode == 200 {
		successCall(body)
	}
	// TO CHECK - What if there was more code in main, and an err has occurred?
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

func successCall(respBody string) {
	fmt.Println(respBody)
}

func client() {
	body, retryAfter, statusCode, _ := getRequest()
	if statusCode == 200 {
		fmt.Println(string(body))
	} else if statusCode == 429 {
		if retryAfter == "a while" {
			fmt.Println("Experiencing high volumes of traffic...retrying again...")
			time.Sleep(time.Duration(4) * time.Second)
			client()
		} else if secs, _ := strconv.Atoi(retryAfter); secs >= 1 && secs < 5 {
			fmt.Println("Retrying in", secs, "seconds...")
			time.Sleep(time.Duration(secs) * time.Second)
			client()
		} else if secs, _ := strconv.Atoi(retryAfter); secs >= 5 {
			fmt.Println("Weather cannot be retrieved :(")
		} else {
			t, _ := http.ParseTime(retryAfter)
			unixTime := t.UTC().Unix()
			currentTime := time.Now().UTC().Unix()
			retrySecs := unixTime - currentTime
			if retrySecs > 5 {
				fmt.Println("Weather cannot be retrieved :(")
			} else {
				fmt.Println("Retrying in", retrySecs, "seconds...")
				time.Sleep(time.Duration(retrySecs) * time.Second)
				client()
			}
		}
	}
}
