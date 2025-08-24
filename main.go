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
	"strings"
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

func (client *CloudreveClient) Post(path string, payload map[string]interface{}) error {
	u := client.base.JoinPath("/api/v4", path)
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(p)
	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respBody))
	return nil
}

func (client *CloudreveClient) Put(path string, payload map[string]interface{}) ([]byte, error) {
	u := client.base.JoinPath("/api/v4", path)
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(p)
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("PUT response:", string(respBody))
	return respBody, nil
}

func (client *CloudreveClient) PutSessionOpenid() (*SessionOpenidResponse, error) {
	params := map[string]interface{}{"hint": nil, "linking": false, "provider": 2}
	content, err := client.Put("/session/openid", params)
	if err != nil {
		return nil, fmt.Errorf("Failed PUT /session/openid: %w", err)
	}

	var resp SessionOpenidResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
			return nil, fmt.Errorf("Failed to parse response: %w", err)
	}
	return &resp, nil
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

	/*
	if len(os.Args) < 3 {
		fmt.Println("Please specify an email")
		os.Exit(1)
	}
	email := os.Args[2]

	if len(os.Args) < 4 {
		fmt.Println("Please specify a password")
		os.Exit(1)
	}
	password := os.Args[3]
	*/
	client := NewCloudreveClient(urlBase)
	fmt.Println("Attempting login at:", urlBase)
	/*url := client.UrlWithQuery("/session/prepare", map[string]string{"email": email})
	err = client.Get(url)
	if err != nil {
		log.Fatal("Failed to prepare session", err)
	}*/

	/*
	body := map[string]interface{}{"email": email, "password": password}
	err = client.Post("/session/token", body)
	*/
	sessionData, err := client.PutSessionOpenid()
	if err != nil {
		log.Fatal("Failed to get session:", err)
	}
	converted := FixJSONEscaping(sessionData.Data)
	fmt.Println("URL:", converted)

}

func FixJSONEscaping(content string) string {
	content = strings.ReplaceAll(content, "\\u003c", "<")
	content = strings.ReplaceAll(content, "\\u003e", ">")
	content = strings.ReplaceAll(content, "\\u0026", "&")
	return content
}

type SessionOpenidResponse struct {
	Code int    `json:"code"`
	Data string  `json:"data"`
	Error *string  `json:"error"`
	Msg string  `json:"msg"`
}
