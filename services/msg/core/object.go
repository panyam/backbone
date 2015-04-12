package core

import (
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

type IObject interface {
	ToDict() map[string]interface{}
	FromDict(map[string]interface{})
}

type Object struct {
	Cls IObject

	/**
	 * Unique system wide ID.
	 */
	Id int64

	/**
	 * When the object was created.
	 */
	Created time.Time

	/**
	 * Status of the user account.
	 * 0 = valid and active
	 * everything else depends on the object type
	 */
	Status int
}
