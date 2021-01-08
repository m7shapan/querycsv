package main

import (
	"fmt"
	"strconv"
	"strings"
)

func exportCMD(words []string) {
	if ok := validateCommandParamters(2, len(words)-1); !ok {
		return
	}

	err := createFileFromTable(words[1], words[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("file exported")
}

func showCMD(words []string) {
	if ok := validateCommandParamters(2, len(words)-1); !ok {
		return
	}

	n, _ := strconv.Atoi(words[2])
	columns, rows, err := tableSelect(words[1], n)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tableRender(columns, rows)
}

func listCMD() {
	err := list()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func loadCMD(words []string) {
	tablename := words[0]
	err := loadFiles(tablename, words[3])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func selectCMD(words []string) {
	err := createTableFromQuery(words[0], strings.Join(words[2:], " "))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func selectAppendCMD(words []string) {
	err := appendToTableFromQuery(words[0], strings.Join(words[2:], " "))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
