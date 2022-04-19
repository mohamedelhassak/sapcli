package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//set custom config file if --config flag passed in cli
func setConfigs(f string) {

	//if --config passed
	if f != "" {

		//check if file exist
		if !utils.IsFileOrDirExists(f) {
			log.Fatalf("[ERROR!...] file not found:  %s", f)
		}
		// get the filepath
		abs, err := filepath.Abs(f)
		if err != nil {
			log.Fatalf("[ERROR!...] reading file from path:  %v", err.Error())
		}

		// get the config name
		base := filepath.Base(abs)

		// get the path
		path := filepath.Dir(abs)

		//
		viper.AddConfigPath(path)
		viper.SetConfigName(strings.Split(base, ".")[0])
		fmt.Println("[INFO!...] Using config file : " + path + "/" + base)

		// Find and read the config file; Handle errors reading the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("[FAILED!...] Failed to read config file: %v", err)
		}
		//str, _ := viper.Get("creds.subscription-id").(string)
		SAP_CLOUD_API_URL = "https://portalrotapi.hana.ondemand.com/v2/subscriptions/" + viper.Get("creds.subscription-id").(string)
		API_TOKEN = viper.Get("creds.api-token").(string)

		//else --config not passed, use defaul config file
	} else {
		readYAMLConfigs(CONF_FILE_NAME, &cfg)
		SAP_CLOUD_API_URL = "https://portalrotapi.hana.ondemand.com/v2/subscriptions/" + cfg.Creds.SubscriptionId
		API_TOKEN = cfg.Creds.ApiToken
	}
}

//read configs from YAML
func readYAMLConfigs(fileName string, cfg *Config) {

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
}

/*//read env vars within the YAML config file
func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("Failed reading env: %s", err)
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

//http get method
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
