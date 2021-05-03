package nbt

import (
	"bufio"
	"reflect"
)

func Unmarshal(reader *bufio.Reader, value reflect.Value) (tagName string, err error) {
	tagType, err := readTagType(reader)
	if err != nil {
		return
	}

	if tagType == tagEnd {
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

func readValue(reader *bufio.Reader, tagType namedTagType, v reflect.Value) error {
	switch tagType {
	case tagByte:
		return readTagByte(reader, v)
	case tagShort:
		return readTagShort(reader, v)
	case tagInt:
		return readTagInt(reader, v)
	case tagLong:
		return readTagLong(reader, v)
	case tagFloat:
		return readTagFloat(reader, v)
	case tagDouble:
		return readTagDouble(reader, v)
	case tagString:
		return readTagString(reader, v)
	case tagList:
		return readTagList(reader, v)
	case tagCompound:
		switch v.Kind() {
		case reflect.Struct:
			return readTagCompoundStruct(reader, v)
		case reflect.Map:
			return readTagCompoundMap(reader, v)
		}
	case tagByteArray:
		return readTagByteArray(reader, v)
	case tagIntArray:
		return readTagIntArray(reader, v)
	case tagLongArray:
		return readTagLongArray(reader, v)
	}
	return nil
}
