package utils

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"strings"
)

func ByteToInt64(data []byte) int64 {
	var value int64
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &value)
	return value
}

func UUIDGen() int64 {
	u1 := uuid.NewV4()
	return ByteToInt64(u1.Bytes())
	// return strings.Replace(u1.String(), "-", "", -1)
}

func CreateTable(db *sql.DB, tableName string, columns []string, tableConstraints ...string) error {
	tableConstraintsString := strings.Join(tableConstraints, ", ")
	columnsString := strings.Join(columns, ", ")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s %s)", tableName, columnsString, tableConstraintsString)
	_, error := db.Exec(query)
	if error != nil {
		log.Printf("Error creating table %s: %s\n", tableName, error)
	}
	return error
}

func InsertRow(db *sql.DB, tableName string, args ...interface{}) error {
	query := fmt.Sprintf("INSERT INTO %s", tableName)
	columnsString := "("
	valuesString := "("
	numArgs := len(args)
	if numArgs%3 != 0 {
		return errors.New("Number of arguments must be a multiple of 3")
	}

	for i := 0; i < numArgs; i += 3 {
		colName := args[i].(string)
		colFormat := args[i+1].(string)
		if i > 0 {
			columnsString = columnsString + ", "
			valuesString = valuesString + ", "
		}
		columnsString = columnsString + " " + colName
		if colFormat == "%d" {
			valuesString = valuesString + fmt.Sprintf(" %d", args[i+2].(int))
		} else if colFormat == "%ld" {
			valuesString = valuesString + fmt.Sprintf(" %d", args[i+2].(int64))
		} else if colFormat == "%s" {
			valuesString = valuesString + fmt.Sprintf(" '%s'", args[i+2].(string))
		} else {
			return errors.New(fmt.Sprintf("Invalid col format: %s at index %d", colFormat, i+1))
		}
	}

	columnsString = columnsString + ")"
	valuesString = valuesString + ")"
	query = query + columnsString + " VALUES " + valuesString
	_, err := db.Exec(query)
	return err
}

func UpdateRows(db *sql.DB, tableName string, whereClause string, args ...interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET", tableName)
	columnsString := "("

	numArgs := len(args)
	if numArgs%3 != 0 {
		return errors.New("Number of arguments must be a multiple of 3")
	}

	for i := 0; i < numArgs; i += 3 {
		colName := args[i].(string)
		colFormat := args[i+1].(string)
		if i > 0 {
			columnsString = columnsString + ", "
		}
		if colFormat == "%d" {
			columnsString += fmt.Sprintf("%s = %d", colName, args[i+2].(int))
		} else if colFormat == "%ld" {
			columnsString += fmt.Sprintf("%s = %d", colName, args[i+2].(int64))
		} else if colFormat == "%s" {
			columnsString += fmt.Sprintf("%s = '%s'", colName, args[i+2].(string))
		} else {
			return errors.New(fmt.Sprintf("Invalid col format: %s at index %d", colFormat, i+1))
		}
	}

	query = query + columnsString
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	_, err := db.Exec(query)
	return err
}

func ClearTable(db *sql.DB, tableName string) error {
	_, error := db.Exec(fmt.Sprintf("DELETE FROM %s CASCADE", tableName))
	if error != nil {
		log.Printf("Error clearing table %s: %s\n", tableName, error)
	}
	return error
}
