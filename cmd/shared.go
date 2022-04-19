package cmd

import (
	"github.com/mohamedelhassak/sapcli/utils"
)

var SAPCLI_VERSION = "SAPCLI v0.1.0"
var WORK_DIR = utils.GetEnvExist("SAPCLI_WORK_DIR", "Environement variable SAPCLI_WORK_DIR not set!")
var CONF_FILE_NAME = WORK_DIR + "/.config.yaml"
var LOGS_DIR = WORK_DIR + "/logs"
var BUILDS_DIR = WORK_DIR + "/builds"

//use default url
var SAP_CLOUD_API_URL = ""
var API_TOKEN = ""
var (
	// One http client that will be used in the whole application
	client = httpClient()

	//init config, cfg instance will be used in all files
	cfg Config
)

func init() {

}
