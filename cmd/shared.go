package cmd

import (
	"github.com/mohamedelhassak/sapcli/utils"
)

var SAPCLI_VERSION = "SAPCLI v0.1.0"
var WORK_DIR = utils.GetEnvExist("SAPCLI_WORK_DIR", "Environement variable SAPCLI_WORK_DIR not set!")
var CONF_FILE_NAME_PATTERN = "*.config.*"

var SLASH_OR_BACKSLASH = ""
var LOGS_DIR = ""
var BUILDS_DIR = ""

//use default url
var SAP_CLOUD_API_URL = ""
var API_TOKEN = ""

// One http client that will be used in the whole application
var client = httpClient()

func init() {
	//return "\" OR "/" according to OS Type
	SLASH_OR_BACKSLASH = utils.SlashOrBackslash()

	//e.g : WORK_DIR/logs/
	LOGS_DIR = WORK_DIR + SLASH_OR_BACKSLASH + "logs" + SLASH_OR_BACKSLASH
	//e.g : WORK_DIR/builds/
	BUILDS_DIR = WORK_DIR + SLASH_OR_BACKSLASH + "builds" + SLASH_OR_BACKSLASH
	//e.g : WORK_DIR/
	WORK_DIR = WORK_DIR + SLASH_OR_BACKSLASH
}
