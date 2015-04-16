package entity

type Field interface {
	Name() string
	Type() string
	Default() interface{}
}

type Entity interface {
	/**
	 * Get the fields of this entity.
	 */
	GetFields() []*Field
	GetField(key string) (interface{}, error)
	SetField(key string, value interface{}) error
}

type Persister interface {
	CreateEntity(entity Entity, override bool) error
	UpdateEntity(entity Entity, override bool) error
	GetEntity(entity Entity)
}
