package core

import (
	"time"
)

type IModel interface {
	GetProperties() []string
	GetProperty(key string) interface{}
	SetProperty(key string, value interface{})
}

type Model struct {
	Cls IObject

	/**
	 * Has the object been loaded?
	 */
	Loaded bool

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
