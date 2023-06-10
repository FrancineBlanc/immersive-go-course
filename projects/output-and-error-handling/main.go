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
		if secs, _ := strconv.Atoi(resp.Header.Get("Retry-After")); secs >= 1 && secs <= 10 {
			fmt.Println("Retrying in", secs, "seconds...")
			time.Sleep(time.Duration(secs) * time.Second)
			client()
		}
	}

}
