package cmd

import (
	"github.com/mohamedelhassak/sapcli/utils"
)

var (
	SAPCLI_VERSION         = "v0.1.0"
	CONF_FILE_NAME_PATTERN = ".config.*"

	WORK_DIR           string
	SLASH_OR_BACKSLASH string
	LOGS_DIR           string
	BUILDS_DIR         string

	//use default url
	SAP_CLOUD_API_URL string
	API_TOKEN         string
)

// One http client that will be used in the whole application
var client = utils.HttpClient()

func init() {
	WORK_DIR = utils.GetEnvExist("SAPCLI_WORK_DIR", "Environement variable SAPCLI_WORK_DIR not set!")

	//return "\" OR "/" according to OS Type
	SLASH_OR_BACKSLASH = utils.SlashOrBackslash()

	//e.g : WORK_DIR/logs/
	LOGS_DIR = WORK_DIR + SLASH_OR_BACKSLASH + "logs" + SLASH_OR_BACKSLASH
	//e.g : WORK_DIR/builds/
	BUILDS_DIR = WORK_DIR + SLASH_OR_BACKSLASH + "builds" + SLASH_OR_BACKSLASH
	//e.g : WORK_DIR/
	WORK_DIR = WORK_DIR + SLASH_OR_BACKSLASH
}
