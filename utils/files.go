package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

//read config file from CSV
func ReadCSVConfig(filePath string) (map[string]string, error) {

	// read csv file
	csvfile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer csvfile.Close()
	csvfile.Seek(0, 0)
	reader := csv.NewReader(csvfile)
	reader.Comma = '|'

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]string)

	for _, record := range rawCSVdata {

		data := strings.Split(record[0], ",")

		dataMap[data[0]] = strings.TrimSpace(data[1])

	}

	return dataMap, err
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

//check & retuen env
func GetEnvExist(key, fallback string) string {
	val, isPresent := os.LookupEnv(key)
	if isPresent == false && val == "" {
		log.Println(fallback)
	}
	return val
}

func SearchFileByPattern(pattern string, dirName string) string {

	//dirName sould ends with "\" or "/" regarding to OS type
	matches, err := filepath.Glob(dirName + pattern)
	if err != nil {
		fmt.Println(err)
	}
	if len(matches) > 1 {
		log.Fatalf("[ERROR!...] More than one config file found")
	}

	fmt.Println("[INFO!...] Config file found: " + matches[0])
	return matches[0]
}
