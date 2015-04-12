package core

type Serializer interface {
	WriteModel(model IModel)
	WriteField(name string, value interface{})
	WriteList(values []interface{})
	WriteMap(values map[string]interface{})
	WriteString(value string)
	WriteInt64(value int64)
	WriteInt32(value int32)
	WriteInt16(value int16)
	WriteInt8(value int8)
	WriteBool(value bool)
	WriteUInt64(value uint64)
	WriteUInt32(value uint32)
	WriteUInt16(value uint16)
	WriteUInt8(value uint8)
	WriteFloat32(value float32)
	WriteFloat64(value float64)
}

type Deserializer interface {
	ReadModel() (IModel, error)
	ReadField(name string, value interface{})
	ReadList() ([]interface{}, error)
	ReadMap() (map[string]interface{}, error)
	ReadString() (string, error)
	ReadInt64() (int64, error)
	ReadInt32() (int32, error)
	ReadInt16() (int16, error)
	ReadInt8() (int8, error)
	ReadBool() (bool, error)
	ReadUInt64() (uint64, error)
	ReadUInt32() (uint32, error)
	ReadUInt16() (uint16, error)
	ReadUInt8() (uint8, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
}
