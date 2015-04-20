package core

import (
	"time"
)

type IModel interface {
	GetProperties() []string
	GetProperty(key string) interface{}
	SetProperty(key string, value interface{})
}
