package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
)

type Metadata struct {
	DB         *sql.DB
	ParentType string
	SG         *ServiceGroup
}

func CreateMetadataTable(db *sql.DB, tableName, parentTableName string, parentType string) {
	parentReferenceColumnName := fmt.Sprintf("%sId", parentType)
	parentReferenceColumnDesc := fmt.Sprintf("%sId bigint NOT NULL REFERENCES %s ON DELETE CASCADE", parentType, parentTableName)
	CreateTable(db, tableName, []string{
		parentReferenceColumnDesc,
		"Key TEXT NOT NULL",
		"Value TEXT DEFAULT ('')",
		"ValueType INT DEFAULT (0)",
		"LastUpdated TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()"},
		fmt.Sprintf(", CONSTRAINT unique_%s_metadata_entry UNIQUE (%s, Key)", parentTableName, parentReferenceColumnName))
}
