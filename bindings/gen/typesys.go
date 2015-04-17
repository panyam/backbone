package gen

type ITypeSystem interface {
}

type TypeSystem struct {
	types       map[string]*Type
	typeCounter int64
}

func NewTypeSystem() *TypeSystem {
	out := TypeSystem{}
	out.types = make(map[string]*Type)
	return &out
}

/**
 * Adds a type to the type system.  If the type already exists then
 * the existing one is returned otherwise a new type is added and returned.
 * Also the type's ID will be set.
 */
func (ts *TypeSystem) AddType(t *Type) (alt *Type) {
	key := t.Package + "." + t.Name
	if value, ok := ts.types[key]; ok {
		return value
	}
	ts.typeCounter++
	t.Id = ts.typeCounter
	ts.types[key] = t
	return t
}

func (ts *TypeSystem) GetType(pkg string, name string) *Type {
	key := pkg + "." + name
	return ts.types[key]
}
