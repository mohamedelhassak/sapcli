package utils

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v2"
)

//return helper message when command unknown
func UnknownCommandMsg(cmd string) string {
	return "[Unknown!...] command not found: Use \"sapcli " + cmd + "--help \" for more information about this command."
}

// PrettyPrint to print JSON struct in a readable way
func PrettyPrintJSON(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		log.Fatalf("[ERROR!...] Couldn't Marshal JSON: %v", err)
	}
	return string(s)
}

// PrettyPrint to print yaml struct in a readable way
func PrettyPrintYAML(i interface{}) string {
	s, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalf("[ERROR!...] Couldn't Marshal yaml: %v", err)
	}
	return string(s)
}
