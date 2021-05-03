package nbt

import (
	"bufio"
	"github.com/junglemc/mc"
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
	switch tagType {
	case tagByte:
		if reflect.TypeOf(value).Kind() == reflect.Bool {
			if value.(bool) {
				return writeByte(writer, 1)
			} else {
				return writeByte(writer, 0)
			}
		}
		return writeByte(writer, value.(byte))
	case tagShort:
		return writeInt16(writer, value.(int16))
	case tagInt:
		return writeInt32(writer, value.(int32))
	case tagLong:
		return writeInt64(writer, value.(int64))
	case tagFloat:
		return writeFloat32(writer, value.(float32))
	case tagDouble:
		return writeFloat64(writer, value.(float64))
	case tagString:
		if reflect.TypeOf(value) == reflect.TypeOf(mc.Identifier{}) {
			return writeString(writer, value.(mc.Identifier).String())
		}
		return writeString(writer, value.(string))
	case tagList:
		return writeList(writer, reflect.ValueOf(value))
	case tagCompound:
		return writeCompound(writer, value)
	case tagByteArray:
		return writeByteSlice(writer, value.([]byte))
	case tagIntArray:
		return writeInt32Slice(writer, value.([]int32))
	case tagLongArray:
		return writeInt64Slice(writer, value.([]int64))
	}
	return nil
}

func typeOf(t reflect.Type) namedTagType {
	if t == reflect.TypeOf(mc.Identifier{}) {
		return tagString
	}

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
