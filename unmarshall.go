package nbt

import (
    "bufio"
    "reflect"
)

func Unmarshall(reader *bufio.Reader, value reflect.Value) (tagName string, err error) {
    tagType, err := readTagType(reader)
    if err != nil {
        return
    }

    if tagType == TagEnd {
        return
    }

    tagName, err = readString(reader)
    if err != nil {
        return
    }

    err = readValue(reader, tagType, value)
    if err != nil {
        return
    }
    return
}

func readValue(reader *bufio.Reader, tagType TagType, vptr reflect.Value) error {
    v := vptr
    if v.Type().Kind() == reflect.Ptr {
        v = v.Elem()
    }
    switch tagType {
    case TagByte:
        return readTagByte(reader, v)
    case TagShort:
        return readTagShort(reader, v)
    case TagInt:
        return readTagInt(reader, v)
    case TagLong:
        return readTagLong(reader, v)
    case TagFloat:
        return readTagFloat(reader, v)
    case TagDouble:
        return readTagDouble(reader, v)
    case TagString:
        return readTagString(reader, v)
    case TagList:
        return readTagList(reader, v)
    case TagCompound:
        switch v.Kind() {
        case reflect.Struct:
            return readTagCompoundStruct(reader, v)
        case reflect.Map:
            return readTagCompoundMap(reader, v)
        }
    case TagByteArray:
        return readTagByteArray(reader, v)
    case TagIntArray:
        return readTagIntArray(reader, v)
    case TagLongArray:
        return readTagLongArray(reader, v)
    }
    return nil
}
