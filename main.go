package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/peterh/liner"
)

var (
	historyFn = filepath.Join(os.TempDir(), ".liner_example_history")
	commands  = []string{"load", "show", "list", "export"}
	db        *sqlx.DB
)

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
}

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	defer func() {
		if f, err := os.Create(historyFn); err == nil {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	for {

		statement, err := line.Prompt("➜ ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				break
			}

			fmt.Println(err.Error())
			continue
		}

		line.AppendHistory(statement)
		words := strings.Split(statement, " ")

		switch {
		case len(words) >= 1 && words[0] == "export":
			exportCMD(words)
			continue
		case len(words) >= 1 && words[0] == "show":
			showCMD(words)
			continue
		case len(words) >= 1 && words[0] == "list":
			listCMD()
			continue
		case len(words) >= 3 && words[1] == "=" && words[2] == "load":
			loadCMD(words)
			continue
		case len(words) >= 3 && words[1] == "=" && words[2] == "select":

		case len(words) >= 3 && words[1] == "+=" && words[2] == "select":
			selectCMD(words)
			continue
		default:
			fmt.Println("no such command supported")
		}
	}
}

func tableRender(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(rows)
	table.Render()
}

func createFileFromTable(tablename, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headerColumns, rows, err := tableSelect(tablename, 0)
	if err != nil {
		return
	}

	// write header columns
	err = writer.Write(headerColumns)
	if err != nil {
		return
	}

	for _, row := range rows {
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func list() (err error) {
	rows, err := db.Query("select name from sqlite_master;")
	if err != nil {
		return
	}
	var tableRows [][]string

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}

		tableRows = append(tableRows, []string{name})
	}

	tableRender([]string{"Tables"}, tableRows)

	return
}

func validateCommandParamters(expected, found int) (ok bool) {
	ok = (expected == found)

	if expected != found {
		fmt.Printf("error: wrong number of command parameters expected %d, found %d \n", expected, found)
		return
	}

	return
}
