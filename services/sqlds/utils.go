package sqlds

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"strings"
	// "errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
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
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s %s)", tableName, strings.Join(columns, ", "), tableConstraintsString)
	_, error := db.Exec(query)
	if error != nil {
		log.Printf("Error creating table %s: %s\n", tableName, error)
	}
	return error
}

func ClearTable(db *sql.DB, tableName string) error {
	_, error := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
	if error != nil {
		log.Printf("Error clearing table %s: %s\n", tableName, error)
	}
	return error
}
