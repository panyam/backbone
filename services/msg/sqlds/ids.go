package sqlds

import (
	"database/sql"
	"errors"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
)

type IDService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const SEQUENCES_TABLE = "sequences"

func NewIDService(db *sql.DB, sg *ServiceGroup) *IDService {
	svc := IDService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *IDService) InitDB() {
	CreateTable(svc.DB, SEQUENCES_TABLE,
		[]string{
			"Domain TEXT PRIMARY KEY",
			"StartVal bigint DEFAULT(1)",
			"Increment int DEFAULT(1)",
		})
}

/**
 * Create a domain of IDs
 */
func (svc *IDService) CreateDomain(domain string, startVal int64, increment int) error {
	svc.DB.QueryRow(fmt.Sprintf("CREATE SEQUENCE %s INCREMENT BY %d MINVALUE %d", domain, increment, startVal))
	InsertRow(svc.DB, SEQUENCES_TABLE,
		"Domain", domain,
		"StartVal", startVal,
		"Increment", increment)
	return nil
}

func (svc *IDService) NextID(domain string) (int64, error) {
	row := svc.DB.QueryRow(fmt.Sprintf("SELECT nextval('%s')", domain))
	var nextVal int64 = 0
	err := row.Scan(&nextVal)
	if err != nil {
		return 0, errors.New("Unable to get sequence")
	}
	return nextVal, nil
}

/**
 * Removes all entries.
 */
func (svc *IDService) RemoveAllSequences() {
	query := fmt.Sprintf("SELECT Domain FROM %s", SEQUENCES_TABLE)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	for rows.Next() {
		var domain string
		rows.Scan(&domain)
		svc.DB.QueryRow("DROP SEQUENCE %s CASCADE", domain)
	}
}
