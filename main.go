package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	files    = "files"
	saveFile = "fileResult/endFile.csv"
)

type CSVData struct {
	TypeFile string
	Count    int64
	Percent  float64
}

type ByCountTypeFile []CSVData

func (a ByCountTypeFile) Len() int           { return len(a) }
func (a ByCountTypeFile) Less(i, j int) bool { return a[i].Count < a[j].Count }
func (a ByCountTypeFile) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	fileList, err := GetListFilesProcess(files)
	if err != nil {
		log.Println(err.Error())
	}

	var numberOfEntries int64 = 0
	mapTypeFiles := make(map[string]int64)
	for _, val := range fileList {
		err := ReadToCSVFile(mapTypeFiles, val, &numberOfEntries)
		if err != nil {
			log.Println(err)
		}
	}

	listPercent := make(ByCountTypeFile, 0)
	for key, val := range mapTypeFiles {
		listPercent = append(listPercent, CSVData{key, val, float64(val) / float64(numberOfEntries)})
	}

	sort.Sort(listPercent)

	WriterCsv(listPercent, saveFile)
}

func GetListFilesProcess(path string) ([]string, error) {
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0, len(lst))

	for _, val := range lst {
		if val.IsDir() {
			fmt.Printf("[%s]\n", val.Name())
		} else {
			files = append(files, path+"/"+val.Name())
		}
	}

	if len(files) == 0 {
		return files, errors.New("no files to process!")
	}
	return files, nil
}

func ReadToCSVFile(mapTypeFiles map[string]int64, file string, count *int64) error {
	filesCSV, err := os.Open(file)
	defer filesCSV.Close()

	if err != nil {
		return err
	}

	lines, err := csv.NewReader(filesCSV).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		a1 := strings.ToLower(line[0])
		a2, _ := strconv.ParseInt(line[1], 10, 64)

		if _, ok := mapTypeFiles[a1]; !ok {
			mapTypeFiles[a1] = a2
		} else {
			mapTypeFiles[a1] += a2
		}
		*count += a2
	}
	return nil
}

func WriterCsv(listFile ByCountTypeFile, pathFile string) {
	records := [][]string{{}}

	for _, val := range listFile {
		records = append(records, []string{val.TypeFile, strconv.Itoa(int(val.Count)), fmt.Sprintf("%f", val.Percent)})
	}

	file, errCreate := os.Create(pathFile)
	if errCreate != nil {
		log.Panic(errCreate)
	}

	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()
}
