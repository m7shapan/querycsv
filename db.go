package main

import (
	"fmt"
	"strconv"
)

func createTableAndInsertRecords(table string, records [][]string) error {
	var tableColumns string
	var questionMarks string
	var queryColumns string

	for i := 0; i < len(records[0]); i++ {
		tableColumns += records[0][i] + " "
		if _, err := strconv.ParseInt(records[1][i], 10, 64); err == nil {
			tableColumns += "int"
		} else {
			tableColumns += "string"
		}

		queryColumns += records[0][i]
		questionMarks += "?"
		if i != (len(records[0]) - 1) {
			tableColumns += ", "
			queryColumns += ", "
			questionMarks += ", "
		}
	}

	// CREATE TABLE tablename ();
	tableCreation := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table, tableColumns)
	_, err := db.Exec(tableCreation)
	if err != nil {
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, queryColumns, questionMarks)
	stmt, err := db.Preparex(insertQuery)
	if err != nil {
		return err
	}

	for i := 1; i < len(records); i++ {

		var columns []interface{}
		for j := 0; j < len(records[i]); j++ {
			if integer, err := strconv.ParseInt(records[i][j], 10, 64); err == nil {
				columns = append(columns, integer)
			} else {
				columns = append(columns, records[i][j])
			}
		}
		if len(columns) == 0 {
			break
		}

		_, err = stmt.Exec(columns...)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTableFromQuery(tablename, query string) (err error) {
	// CREATE TABLE new_table_name AS
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s AS %s", tablename, query))
	if err != nil {
		return
	}

	return
}

func appendToTableFromQuery(tablename, query string) (err error) {
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s %s", tablename, query))
	if err != nil {
		return
	}

	return
}

func tableSelect(table string, limit int) (headerColumns []string, rows [][]string, err error) {
	sql := fmt.Sprintf("SELECT * FROM %s", table)
	if limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d", sql, limit)
	}

	dbRows, err := db.Queryx(sql)
	if err != nil {
		return
	}

	headerColumns, err = dbRows.Columns()
	if err != nil {
		return
	}

	for dbRows.Next() {
		var result []interface{}
		result, err = dbRows.SliceScan()
		if err != nil {
			return
		}

		var stringColumns []string
		for i := 0; i < len(result); i++ {
			stringColumns = append(stringColumns, fmt.Sprint(result[i]))
		}

		rows = append(rows, stringColumns)
	}

	return
}
