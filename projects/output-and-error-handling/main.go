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
	_, err, _, _ := getRequest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		client()
	}
}

func getRequest() (string, error, int, string) {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		return "", fmt.Errorf("the weather cannot be retrieved right now"), 0, ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read weather data at present"), 0, ""
	}
	return string(body), nil, resp.StatusCode, resp.Header.Get("Retry-After")
}

func client() {
	body, _, statusCode, retryAfter := getRequest()
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
			fmt.Println("Retrying in", retrySecs, "seconds...")
			time.Sleep(time.Duration(retrySecs) * time.Second)
			client()
		}
	}

}
