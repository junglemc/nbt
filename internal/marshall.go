package internal

import (
	"errors"
	"reflect"
)

func (c *Codec) Marshall(tagName string, tagType TagType, v reflect.Value) (err error) {
	if err = c.writeHeader(tagName, tagType); err != nil {
		return err
	}
	if err = c.writeValue(tagType, v); err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeHeader(tagName string, tagType TagType) (err error) {
	if err = c.writeTagType(tagType); err != nil {
		return err
	}

	if err = c.writeString(tagName); err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeValue(tagType TagType, v reflect.Value) error {
	switch tagType {
	case TagByte:
		return c.writeByte(byte(v.Uint()))
	case TagShort:
		return c.writeInt16(int16(v.Int()))
	case TagInt:
		return c.writeInt32(int32(v.Int()))
	case TagLong:
		return c.writeInt64(v.Int())
	case TagFloat:
		return c.writeFloat32(float32(v.Float()))
	case TagDouble:
		return c.writeFloat64(v.Float())
	case TagString:
		return c.writeString(v.String())
	case TagList:
		nestedTagType := c.TypeOf(v.Type().Elem())
		if v.Len() <= 0 {
			nestedTagType = TagEnd // Mimic notchian behavior
		}

		var err error
		if err = c.writeTagType(nestedTagType); err != nil {
			return err
		}

		n := v.Len()
		if err = c.writeInt32(int32(n)); err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			val := v.Index(i)
			err = c.writeValue(nestedTagType, val)
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

			nestedTagType := c.TypeOf(f.Type)
			if f.Tag.Get("nbt_type") == "list" {
				nestedTagType = TagList
			}

			err := c.Marshall(nestedTagName, nestedTagType, v.Field(i))
			if err != nil {
				return err
			}
		}
		return c.writeTagType(TagEnd)
	case TagByteArray:
		return c.writeByteSlice(v.Bytes())
	case TagIntArray, TagLongArray:
		n := v.Len()
		if err := c.writeInt32(int32(n)); err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			val := v.Index(i).Int()

			var err error
			if tagType == TagIntArray {
				err = c.writeInt32(int32(val))
			} else if tagType == TagLongArray {
				err = c.writeInt64(val)
			}
			if err != nil {
				return err
			}
		}
	}
	return errors.New("unknown tag type")
}

func (c *Codec) TypeOf(vk reflect.Type) TagType {
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
