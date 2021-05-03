package nbt

import (
	"bufio"
	"reflect"
)

func Marshal(writer *bufio.Writer, tagName string, value interface{}) (err error) {
	var tagType namedTagType
	if value == nil {
		tagType = tagCompound
	} else {
		tagType = typeOf(reflect.TypeOf(value))
	}

	err = writeTagType(writer, tagType)
	if err != nil {
		return
	}

	err = writeString(writer, tagName)
	if err != nil {
		return
	}

	err = writeValue(writer, tagType, value)
	if err != nil {
		return
	}
	return
}

func writeValue(writer *bufio.Writer, tagType namedTagType, value interface{}) error {
	v := reflect.ValueOf(value)

	switch tagType {
	case tagByte:
		if reflect.TypeOf(value).Kind() == reflect.Bool {
			if v.Bool() {
				return writeByte(writer, 1)
			} else {
				return writeByte(writer, 0)
			}
		}
		return writeByte(writer, byte(v.Uint()))
	case tagShort:
		return writeInt16(writer, int16(v.Int()))
	case tagInt:
		return writeInt32(writer, int32(v.Int()))
	case tagLong:
		return writeInt64(writer, v.Int())
	case tagFloat:
		return writeFloat32(writer, float32(v.Float()))
	case tagDouble:
		return writeFloat64(writer, v.Float())
	case tagString:
		return writeString(writer, v.String())
	case tagList:
		return writeList(writer, v)
	case tagCompound:
		return writeCompound(writer, value)
	case tagByteArray:
		return writeByteSlice(writer, v.Bytes())
	case tagIntArray:
		return writeInt32Slice(writer, v)
	case tagLongArray:
		return writeInt64Slice(writer, v)
	}
	return nil
}

func typeOf(t reflect.Type) namedTagType {
	switch t.Kind() {
	case reflect.Uint8, reflect.Bool:
		return tagByte
	case reflect.Int16, reflect.Uint16:
		return tagShort
	case reflect.Int32, reflect.Uint32:
		return tagInt
	case reflect.Float32:
		return tagFloat
	case reflect.Int64, reflect.Uint64:
		return tagLong
	case reflect.Float64:
		return tagDouble
	case reflect.String:
		return tagString
	case reflect.Struct, reflect.Interface, reflect.Map:
		return tagCompound
	case reflect.Array, reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return tagByteArray
		case reflect.Int32:
			return tagIntArray
		case reflect.Int64:
			return tagLongArray
		default:
			return tagList
		}
	default:
		return tagNone
	}
}
