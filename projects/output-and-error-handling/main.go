package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	}
}
