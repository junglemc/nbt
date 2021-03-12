package internal

import (
	"errors"
	"fmt"
	"reflect"
)

func (c *Codec) Unmarshall(v *reflect.Value) (string, error) {
	tagType, tagName, err := c.readHeader()
	if err != nil {
		return "", err
	}

	e := v.Elem()

	if err = c.readValue(tagType, &e); err != nil {
		return "", err
	}

	return tagName, nil
}

func (c *Codec) readHeader() (TagType, string, error) {
	tagType, err := c.readTagType()
	if err != nil {
		return TagNone, "", err
	}

	if tagType == TagEnd {
		return TagEnd, "", nil
	}

	tagName, err := c.readString()
	if err != nil {
		return TagNone, "", err
	}
	return tagType, tagName, nil
}

func (c *Codec) readValue(tagType TagType, v *reflect.Value) error {
	switch tagType {
	case TagByte:
		value, err := c.readByte()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(value))
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagByte as " + kind.String())
		}
	case TagShort:
		value, err := c.readInt16()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(value))
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(value))
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagShort as " + kind.String())
		}
	case TagInt:
		value, err := c.readInt32()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Int, reflect.Int32, reflect.Int64:
			v.SetInt(int64(value))
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(value))
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagInt as " + kind.String())
		}
	case TagLong:
		value, err := c.readInt64()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Int, reflect.Int64:
			v.SetInt(value)
		case reflect.Uint, reflect.Uint64:
			v.SetUint(uint64(value))
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagLong as " + kind.String())
		}
	case TagFloat:
		value, err := c.readFloat32()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Float32:
			v.SetFloat(float64(value))
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagFloat as " + kind.String())
		}
	case TagDouble:
		value, err := c.readFloat64()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.Float32, reflect.Float64:
			v.SetFloat(value)
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagDouble as " + kind.String())
		}
	case TagString:
		value, err := c.readString()
		if err != nil {
			return err
		}

		switch kind := v.Kind(); kind {
		case reflect.String:
			v.SetString(value)
		case reflect.Interface:
			v.Set(reflect.ValueOf(value))
		default:
			return errors.New("cannot parse TagString as " + kind.String())
		}
	case TagList:
		listType, err := c.readTagType()
		if err != nil {
			return err
		}

		length, err := c.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			length = 0
		}

		var list reflect.Value
		kind := v.Kind()

		switch kind {
		case reflect.Interface:
			list = reflect.ValueOf(make([]interface{}, length))
		case reflect.Slice:
			list = reflect.MakeSlice(v.Type(), int(length), int(length))
		case reflect.Array:
			if arrayLength := v.Len(); arrayLength < int(length) {
				return fmt.Errorf("size mismatch in TagList: want=%d, available=%d", arrayLength, length)
			}
			list = *v
		}

		for i := 0; i < int(length); i++ {
			ind := list.Index(i)
			err = c.readValue(listType, &ind)
			if err != nil {
				return err
			}
		}

		if kind != reflect.Array {
			v.Set(list)
		}
		return nil
	case TagCompound:
		switch kind := v.Kind(); kind {
		case reflect.Struct:
			t := v.Type()
			indices := make(map[string]int)

			n := v.NumField()
			for i := 0; i < n; i++ {
				f := t.Field(i)
				tag := f.Tag.Get("nbt")
				if tag == "-" {
					continue
				}

				if tag != "" {
					indices[tag] = i
				} else {
					indices[f.Name] = i
				}
			}

			for  {
				cmpTagType, cmpTagName, err := c.readHeader()

				if cmpTagType == TagEnd {
					return nil
				}

				if err != nil {
					return err
				}

				index, ok := indices[cmpTagName]
				if !ok {
					return errors.New("no name index found")
				}

				f := v.Field(index)
				if err = c.readValue(cmpTagType, &f); err != nil {
					return err
				}
			}
		case reflect.Map:
			if v.Type().Key().Kind() != reflect.String {
				return errors.New("map key should be of type string")
			}

			if v.IsNil() {
				v.Set(reflect.MakeMap(v.Type()))
			}

			for {
				cmpTagType, cmpTagName, err := c.readHeader()

				if cmpTagType == TagEnd {
					return nil
				}

				if err != nil {
					return err
				}
				val := reflect.New(v.Type().Elem())
				if err = c.readValue(cmpTagType, &val); err != nil {
					return err
				}
				v.SetMapIndex(reflect.ValueOf(cmpTagName), v.Elem())
			}
		}
	case TagByteArray:
		b, err := c.readByteSlice()
		if err != nil {
			return err
		}

		switch t := v.Type(); {
		case t == reflect.TypeOf(b):
			v.SetBytes(b)
		case v.Kind() == reflect.Interface:
			v.Set(reflect.ValueOf(b))
		}
	case TagIntArray:
		length, err := c.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			length = 0
		}

		t := v.Type()
		if t.Kind() == reflect.Interface {
			t = reflect.TypeOf([]int32{})
		} else if t.Kind() != reflect.Slice {
			return errors.New("slice required")
		} else if ek := t.Elem().Kind(); ek != reflect.Int && ek != reflect.Int32 {
			return errors.New("slice of int or int32 type required")
		}

		b := reflect.MakeSlice(t, int(length), int(length))
		for i:=0; i<int(length); i++ {
			val, err := c.readInt32()
			if err != nil {
				return err
			}
			b.Index(i).SetInt(int64(val))
		}
		v.Set(b)
	case TagLongArray:
		length, err := c.readInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			length = 0
		}

		t := v.Type()
		if t.Kind() == reflect.Interface {
			t = reflect.TypeOf([]int64{})
		} else if t.Kind() != reflect.Slice {
			return errors.New("slice required")
		} else if ek := t.Elem().Kind(); ek != reflect.Int && ek != reflect.Int64 {
			return errors.New("slice of int or int64 type required")
		}

		b := reflect.MakeSlice(t, int(length), int(length))
		for i:=0; i<int(length); i++ {
			val, err := c.readInt64()
			if err != nil {
				return err
			}
			b.Index(i).SetInt(val)
		}
		v.Set(b)
	}
	return nil
}
