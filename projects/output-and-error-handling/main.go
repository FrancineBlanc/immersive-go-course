package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Response Status: ", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

// TODO: write second function that returns the error from the GET req, have main call this func, and if error == nil main will finish, else can use log fatal if error can't be reversed.
