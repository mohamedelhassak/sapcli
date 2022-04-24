package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
)

//create directory
func CreateDir(dirName string) {
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

//write into file
func WriteFile(dirName string, fileName string, data string) (int, error) {

	//dirName sould ends with "\" or "/" regarding to OS type
	if !IsFileOrDirExists(dirName) {
		CreateDir(dirName)
	}

	file, err := os.Create(dirName + fileName)
	if err != nil {
		log.Fatalf("[FAILED!...] Failed creating file: %s", err)
	}
	n, err := file.WriteString(data)

	file.Close()

	return n, err
}

//download & save .zip file
func DownloadZipFile(dirName string, zipFileName string, body []byte) error {
	buf := new(bytes.Buffer)

	if _, err := buf.Write(body); err != nil {
		log.Fatal(err)
	}

	//dirName sould ends with "\" or "/" regarding to OS type
	if !IsFileOrDirExists(dirName) {
		CreateDir(dirName)
	}

	file, err := os.Create(dirName + zipFileName)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(file, buf); err != nil {
		log.Fatal(err)
	}
	file.Close()

	return err
}

//returns whether the given file or directory exists
func IsFileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//check & retuen env if exist
func GetEnvExist(key, fallback string) string {
	val, isPresent := os.LookupEnv(key)
	if !isPresent && val == "" {
		log.Fatalln(fallback)
	}
	return val
}

func SearchFileByPattern(pattern string, dirName string) string {

	//dirName sould ends with "\" or "/" regarding to OS type
	matches, err := filepath.Glob(dirName + pattern)
	if err != nil {
		log.Fatalf("[ERROR!...] %v", err)
	}
	if len(matches) > 1 {
		log.Fatalf("[ERROR!...] %d config file(s) found, accepts 1", len(matches))
	} else if len(matches) == 0 {
		log.Fatalf("[ERROR!...] No config file found, accepts 1 but found 0")
	}
	return matches[0]
}
