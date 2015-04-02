package core

import (
	"time"
)

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
	 * everything else = invalid
	 */
	Status int
}
