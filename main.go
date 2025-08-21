package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type CloudreveClient struct {
	base url.URL
	httpClient *http.Client
}

func NewCloudreveClient(baseUrl *url.URL) *CloudreveClient {
	return &CloudreveClient{
		base: *baseUrl,
		httpClient: &http.Client{},
	}
}
func (client *CloudreveClient) UrlWithQuery(path string, query map[string]string) *url.URL {
	url := client.base.JoinPath("/api/v4", path)

	q := url.Query()
	for k, v := range query {
			q.Add(k, v)
	}
	url.RawQuery = q.Encode()
	return url
}

func (client *CloudreveClient) Get(u *url.URL) error {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func (client *CloudreveClient) GetWithBody(path string, payload map[string]interface{}) error {
	url := client.base.JoinPath("/api/v4", path)
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(p)
	req, err := http.NewRequest(http.MethodGet, url.String(), body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	log.Println("Status:", resp.Status)
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(respBody))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify the root URL")
		os.Exit(1)
	}
	urlBase, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("Failed to parse base URL", err)
	}

	if len(os.Args) < 3 {
		fmt.Println("Please specify an email")
		os.Exit(1)
	}
	email := os.Args[2]

	client := NewCloudreveClient(urlBase)
	fmt.Println("Attempting login at:", urlBase)
	url := client.UrlWithQuery("/session/prepare", map[string]string{"email": email})
	err = client.Get(url)
	if err != nil {
		log.Fatal("Failed to prepare session", err)
	}
}
