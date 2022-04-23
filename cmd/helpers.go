package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
)

//read configs from YAML (not used)
/*func readYAMLConfigs(fileName string, cfg *Config) {
	fmt.Println("[INFO!...] Using default config file :" + fileName)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("[FAILED!...] Failed reading configs: %s", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatalf("[ERROR!...] Couldn't parse yaml configs: %s", err)
	}
}*/

//return http client
func httpClient() *http.Client {
	client := &http.Client{}
	return client
}

//http get method
func httpGet(client *http.Client, url string) (body []byte) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

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

//http get method (not used)
func httpGetV2(client *http.Client, url string, target interface{}) error {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("[INFO!...]HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode == http.StatusOK {
		//We Read the response body on the line below.
		return json.NewDecoder(resp.Body).Decode(target)

	} else {
		return errors.New("[FAILED!...] Request Failed " + fmt.Sprint(resp.StatusCode))
	}

}

//http post
func httpPost(client *http.Client, url string, reqBody []byte) []byte {

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

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

//validate command args
func isValidArgs(cmd *cobra.Command, ValidArgs []string, args []string) error {
	if len(args) > 1 {
		return errors.New("Requires at least one arg")
	}

	if utils.Contains(ValidArgs, args[0]) {
		return nil
	} else {
		return errors.New("Invalid argument: " + args[0])
	}
}

//return error if more than one arg passed
func isOneAndOnlyValidArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("\"" + cmd.CommandPath() + "\" Requires exactly 1 arg.")
	}
	return cobra.OnlyValidArgs(cmd, args)
}
