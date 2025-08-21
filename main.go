package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify the root URL")
		os.Exit(1)
	}
	url := os.Args[1]
	fmt.Println("Attempting login at:", url)
	//resp, err := http.Get(URL + "/session/prepare")
	resp, err := http.Get(url + "/api/v4/site/ping")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status:", resp.Status)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
