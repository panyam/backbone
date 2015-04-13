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
func (svc *IDService) CreateDomain(request *CreateDomainRequest) error {
	svc.DB.QueryRow(fmt.Sprintf("CREATE SEQUENCE %s INCREMENT BY %d MINVALUE %d",
		request.Domain, request.Increment, request.StartValue))
	InsertRow(svc.DB, SEQUENCES_TABLE,
		"Domain", request.Domain,
		"StartVal", request.StartValue,
		"Increment", request.Increment)
	return nil
}

func (svc *IDService) NextID(request *NextIDRequest) (int64, error) {
	row := svc.DB.QueryRow(fmt.Sprintf("SELECT nextval('%s')", request.Domain))
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
func (svc *IDService) RemoveAllDomains(request *Request) {
	query := fmt.Sprintf("SELECT Domain FROM %s", SEQUENCES_TABLE)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	for rows.Next() {
		var domain string
		rows.Scan(&domain)
		svc.DB.QueryRow(fmt.Sprintf("ALTER SEQUENCE %s START 1", domain))
	}
}
