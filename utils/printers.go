package utils

import (
	"encoding/json"
	"log"
	"runtime"

	"gopkg.in/yaml.v2"
)

// PrettyPrint to print JSON struct in a readable way
func PrettyFormatJSON(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		log.Fatalf("[ERROR!...] Couldn't Marshal JSON: %v", err)
	}
	return string(s)
}

// PrettyPrint to print yaml struct in a readable way
func PrettyFormatYAML(i interface{}) string {
	s, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalf("[ERROR!...] Couldn't Marshal yaml: %v", err)
	}
	return string(s)
}

//PrettyPrint returns string according to the extention eg:json or yaml
func PrettyFormat(i interface{}, extension string) string {
	switch extension {
	case ".json":
		return PrettyFormatJSON(i)
	case ".yaml", ".yml":
		return PrettyFormatYAML(i)
	default:
		return "No valid extention"
	}
}

//SlashOrBackslash returns "\" if OS is windows, "/" otherwise
func SlashOrBackslash() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "\\"
	case "darwin":
		return "/"
	case "linux":
		return "/"
	default:
		return "/"
	}
}
