package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	client()
}

func client() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println(string(body))
	} else if resp.StatusCode == 429 {
		if resp.Header.Get("Retry-After") == "a while" {
			fmt.Println("Experiencing high amounts of traffic...retrying again...")
			time.Sleep(time.Duration(4) * time.Second)
			client()
		} else if secs, _ := strconv.Atoi(resp.Header.Get("Retry-After")); secs >= 1 && secs < 5 {
			fmt.Println("Retrying in", secs, "seconds...")
			time.Sleep(time.Duration(secs) * time.Second)
			client()
		} else if secs, _ := strconv.Atoi(resp.Header.Get("Retry-After")); secs >= 5 {
			fmt.Println("Weather cannot be retrieved :(")
		} else {
			t, _ := http.ParseTime(resp.Header.Get("Retry-After"))
			unixTime := t.UTC().Unix()
			currentTime := time.Now().UTC().Unix()
			retrySecs := unixTime - currentTime
			fmt.Println("Retrying in", retrySecs, "seconds...")
			time.Sleep(time.Duration(retrySecs) * time.Second)
			client()
		}
	}

}
