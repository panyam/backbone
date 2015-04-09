package utils

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"reflect"
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

func formatSqlValue(value interface{}) string {
	if value == nil {
		return "NULL"
	}
	v := reflect.ValueOf(value)
	// t := v.Type()
	kind := v.Kind()
	if kind == reflect.Bool {
		if value.(bool) == true {
			return "true"
		}
		return "false"
	} else if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
		return fmt.Sprintf("%d", value)
	} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return fmt.Sprintf("%d", value)
	} else if kind == reflect.String {
		// TODO: Format this to prevent injection
		return fmt.Sprintf("'%s'", value)
	}
	return "'<unknown_type>'"
}

func InsertRow(db *sql.DB, tableName string, args ...interface{}) error {
	query := fmt.Sprintf("INSERT INTO %s", tableName)
	columnsString := "("
	valuesString := "("
	numArgs := len(args)
	if numArgs%2 != 0 {
		return errors.New("Number of arguments must be a multiple of 2")
	}

	for i := 0; i < numArgs; i += 2 {
		colName := args[i].(string)
		if i > 0 {
			columnsString += ", "
			valuesString += ", "
		}
		columnsString += colName
		valuesString += formatSqlValue(args[i+1])
	}

	columnsString = columnsString + ")"
	valuesString = valuesString + ")"
	query = query + columnsString + " VALUES " + valuesString
	_, err := db.Exec(query)
	return err
}

func UpdateRows(db *sql.DB, tableName string, whereClause string, args ...interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	columnsString := ""

	numArgs := len(args)
	if numArgs%2 != 0 {
		return errors.New("Number of arguments must be a multiple of 2")
	}

	for i := 0; i < numArgs; i += 2 {
		colName := args[i].(string)
		if i > 0 {
			columnsString = columnsString + ", "
		}
		columnsString += fmt.Sprintf("%s = %s", colName, formatSqlValue(args[i+1]))
	}

	query = query + columnsString
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("No rows found")
	}
	return err
}

func DeleteById(db *sql.DB, tableName string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE Id = %d", tableName, id)
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
