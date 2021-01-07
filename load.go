package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func loadFiles(tablename, filespath string) (err error) {
	if filespath == "." || filespath == "./" {
		filespath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	p, err := os.Stat(filespath)
	if err != nil {
		return err
	}

	var files []string
	if p.Mode().IsRegular() {
		files = append(files, filespath)
	} else {
		err := filepath.Walk(filespath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && path[len(path)-4:] == ".csv" {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if len(files) == 0 {
		return errors.New("file not found")
	}

	for _, file := range files {
		fmt.Println("- loading file", file)
		records, err := loadCSV(file)
		if err != nil {
			return err
		}

		err = createTableAndInsertRecords(tablename, records)
		if err != nil {
			return err
		}
	}

	columns, rows, err := tableSelect(tablename, 100)
	if err != nil {
		return err
	}

	tableRender(columns, rows)
	return nil
}

func loadCSV(filename string) (records [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	r := csv.NewReader(file)
	return r.ReadAll()
}
