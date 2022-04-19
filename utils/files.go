package utils

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
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

	if !IsFileOrDirExists(dirName) {
		CreateDir(dirName)
	}

	file, err := os.Create(dirName + "/" + fileName)
	if err != nil {
		log.Fatalf("[FAILED!...] Failed creating file: %s", err)
	}
	n, err := file.WriteString(data)

	file.Close()

	return n, err
}

//set config into CSV
func SetCSVConfig(fileName string, data [][]string) error {
	csvFile, err := os.Create(fileName)

	if err != nil {
		log.Fatalf("[ERROR!...] Failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	for _, confRow := range data {
		err = csvwriter.Write(confRow)
	}

	csvwriter.Flush()
	csvFile.Close()

	return err
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

	if !IsFileOrDirExists(dirName) {
		CreateDir(dirName)
	}

	file, err := os.Create(dirName + "/" + zipFileName)
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