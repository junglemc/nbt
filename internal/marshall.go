package internal

import (
	"bufio"
	"errors"
	"reflect"
)

func Marshall(writer *bufio.Writer, tagName string, tagType TagType, v reflect.Value) (err error) {
	if err = writeHeader(writer, tagName, tagType); err != nil {
		return err
	}
	if err = writeValue(writer, tagType, v); err != nil {
		return err
	}
	return nil
}

func writeHeader(writer *bufio.Writer, tagName string, tagType TagType) (err error) {
	if err = writeTagType(writer, tagType); err != nil {
		return err
	}

	if err = writeString(writer, tagName); err != nil {
		return err
	}
	return nil
}

func writeValue(writer *bufio.Writer, tagType TagType, v reflect.Value) error {
	switch tagType {
	case TagByte:
		return writeByte(writer, byte(v.Uint()))
	case TagShort:
		return writeInt16(writer, int16(v.Int()))
	case TagInt:
		return writeInt32(writer, int32(v.Int()))
	case TagLong:
		return writeInt64(writer, v.Int())
	case TagFloat:
		return writeFloat32(writer, float32(v.Float()))
	case TagDouble:
		return writeFloat64(writer, v.Float())
	case TagString:
		return writeString(writer, v.String())
	case TagList:
		nestedTagType := TypeOf(v.Type().Elem())
		if v.Len() <= 0 {
			nestedTagType = TagEnd // Mimic notchian behavior
		}

		var err error
		if err = writeTagType(writer, nestedTagType); err != nil {
			return err
		}

		n := v.Len()
		if err = writeInt32(writer, int32(n)); err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			val := v.Index(i)
			err = writeValue(writer, nestedTagType, val)
			if err != nil {
				return err
			}
		}
		return nil
	case TagCompound:
		if v.Kind() == reflect.Interface {
			v = reflect.ValueOf(v.Interface())
		}

		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			f := v.Type().Field(i)

			nestedTagName := f.Tag.Get("nbt")

			// Take the field name if unspecified
			if nestedTagName == "" {
				nestedTagName = f.Name
			}

			nestedTagType := TypeOf(f.Type)
			if f.Tag.Get("nbt_type") == "list" {
				nestedTagType = TagList
			}

			err := Marshall(writer, nestedTagName, nestedTagType, v.Field(i))
			if err != nil {
				return err
			}
		}
		return writeTagType(writer, TagEnd)
	case TagByteArray:
		return writeByteSlice(writer, v.Bytes())
	case TagIntArray, TagLongArray:
		n := v.Len()
		if err := writeInt32(writer, int32(n)); err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			val := v.Index(i).Int()

			var err error
			if tagType == TagIntArray {
				err = writeInt32(writer, int32(val))
			} else if tagType == TagLongArray {
				err = writeInt64(writer, val)
			}
			if err != nil {
				return err
			}
		}
	}
	return errors.New("unknown tag type")
}

func TypeOf(vk reflect.Type) TagType {
	switch vk.Kind() {
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
	case reflect.Struct, reflect.Interface:
		return TagCompound
	case reflect.Array, reflect.Slice:
		switch vk.Elem().Kind() {
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
