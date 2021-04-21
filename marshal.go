package nbt

import (
    "bufio"
    "reflect"
)

func Marshal(writer *bufio.Writer, tagName string, value interface{}) (err error) {
    tagType := TypeOf(reflect.TypeOf(value))

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

func writeValue(writer *bufio.Writer, tagType TagType, value interface{}) error {
    switch tagType {
    case TagByte:
        return writeByte(writer, value.(byte))
    case TagShort:
        return writeInt16(writer, value.(int16))
    case TagInt:
        return writeInt32(writer, value.(int32))
    case TagLong:
        return writeInt64(writer, value.(int64))
    case TagFloat:
        return writeFloat32(writer, value.(float32))
    case TagDouble:
        return writeFloat64(writer, value.(float64))
    case TagString:
        return writeString(writer, value.(string))
    case TagList:
        return writeList(writer, reflect.ValueOf(value))
    case TagCompound:
        return writeCompound(writer, reflect.ValueOf(value))
    case TagByteArray:
        return writeByteSlice(writer, value.([]byte))
    case TagIntArray:
        return writeInt32Slice(writer, value.([]int32))
    case TagLongArray:
        return writeInt64Slice(writer, value.([]int64))
    }
    return nil
}

func TypeOf(t reflect.Type) TagType {
    switch t.Kind() {
    case reflect.Uint8:
        return TagByte
    case reflect.Int16, reflect.Uint16:
        return TagShort
    case reflect.Int32, reflect.Uint32:
        return TagInt
    case reflect.Float32:
        return TagFloat
    case reflect.Int64, reflect.Uint64:
        return TagLong
    case reflect.Float64:
        return TagDouble
    case reflect.String:
        return TagString
    case reflect.Struct, reflect.Interface, reflect.Map:
        return TagCompound
    case reflect.Array, reflect.Slice:
        switch t.Elem().Kind() {
        case reflect.Uint8:
            return TagByteArray
        case reflect.Int32:
            return TagIntArray
        case reflect.Int64:
            return TagLongArray
        default:
            return TagList
        }
    default:
        return TagNone
    }
}
