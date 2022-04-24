package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//return http client
func HttpClient() *http.Client {
	client := &http.Client{}
	return client
}

//send http get & return body response
func HttpGet(client *http.Client, url string, apiToken string) (body []byte) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("[INFO!...] HTTP GET Status Code:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode == http.StatusOK {
		//We Read the response body on the line below.
		body, err = ioutil.ReadAll(resp.Body) // response body is []byte
		if err != nil {
			log.Fatalf("[ERROR!...] Couldn't parse response body. %+v", err)
		}
		return body
	} else {
		fmt.Println("[FAILED!...]Request Failed " + fmt.Sprint(resp.StatusCode))
	}
	return
}

//send http post & return body response
func HttpPost(client *http.Client, url string, apiToken string, reqBody []byte) []byte {

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("[INFO!...] HTTP POST Status Code:", resp.StatusCode, http.StatusText(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR!...]Couldn't parse response body. %+v", err)
	}
	return body
}
