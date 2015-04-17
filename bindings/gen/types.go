package gen

const (
	NullType = iota
	LazyType
	BasicType
	AliasType
	ReferenceType
	TupleType
	RecordType
	FunctionType
	ArrayType
	MapType
)

type Type struct {
	// One of the above type classes
	TypeClass int

	// A unique ID assigned to each type
	TypeId int64

	// Package where the type resides
	Package string

	// Name of the type - can be empty for unnamed types (eg functions)
	Name string

	TypeData interface{}
}

type AliasTypeData struct {
	// Type this is an alias/typedef for
	AliasFor *Type
}

type ReferenceTypeData struct {
	// The target type this is a reference to
	TargetType *Type
}

type MapTypeData struct {
	// The target type this is an array of
	KeyType   *Type
	ValueType *Type
}

type ArrayTypeData struct {
	// The target type this is an array of
	TargetType *Type
}

type TupleTypeData struct {
	SubTypes []*Type
}

type RecordTypeData struct {
	// Type of each member in the struct
	InheritedTypes []*Type
	FieldTypes     []*Type
	FieldNames     []string
}

type FunctionTypeData struct {
	// Types of the input parameters
	InputTypes []*Type

	// Types of the output parameters
	OutputTypes []*Type

	// Types of possible exceptions that can be thrown (not supported in all languages)
	ExceptionTypes []*Type
}

func (t *Type) Signature() string {
	switch t.TypeClass {
	case NullType:
		return ""
	case LazyType:
		return t.Name
	case BasicType:
		return t.Name
	case AliasType:
		return t.Name
	case ReferenceType:
		return "*" + t.TypeData.(ReferenceTypeData).TargetType.Signature()
	case RecordType:
		return t.Name
	case TupleType:
		out := "("
		for index, childType := range t.TypeData.(TupleTypeData).SubTypes {
			if index > 0 {
				out += ", "
			}
			out += childType.Signature()
		}
		return out
	case FunctionType:
		funcTypeData := t.TypeData.(FunctionTypeData)
		out := "func"
		out += TypeListSignature(funcTypeData.InputTypes)
		if funcTypeData.OutputTypes != nil {
			out += TypeListSignature(funcTypeData.OutputTypes)
		}
		if funcTypeData.ExceptionTypes != nil {
			out += " throws" + TypeListSignature(funcTypeData.ExceptionTypes)
		}
		return out
	case ArrayType:
		return "[]" + t.TypeData.(ArrayTypeData).TargetType.Signature()
	case MapType:
		mapTypeData := t.TypeData.(MapTypeData)
		return "map[" + mapTypeData.KeyType.Signature() + "]" + mapTypeData.ValueType.Signature()
	}
	return ""
}

func TypeListSignature(types []*Type) string {
	out := " ("
	if types != nil {
		for index, inType := range types {
			if index > 0 {
				out += ","
			}
			out += inType.Signature()
		}
	}
	out += ")"
	return out
}
