package nbt

import (
    "bufio"
    "errors"
    "fmt"
    "reflect"
)

func readTagByte(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readByte(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        v.SetInt(int64(value))
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        v.SetUint(uint64(value))
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagByte as " + kind.String())
    }
    return
}

func readTagShort(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readInt16(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
        v.SetInt(int64(value))
    case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        v.SetUint(uint64(value))
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagShort as " + kind.String())
    }
    return
}

func readTagInt(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readInt32(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Int, reflect.Int32, reflect.Int64:
        v.SetInt(int64(value))
    case reflect.Uint, reflect.Uint32, reflect.Uint64:
        v.SetUint(uint64(value))
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagInt as " + kind.String())
    }
    return
}

func readTagLong(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readInt64(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Int, reflect.Int64:
        v.SetInt(value)
    case reflect.Uint, reflect.Uint64:
        v.SetUint(uint64(value))
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagLong as " + kind.String())
    }
    return
}

func readTagFloat(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readFloat32(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Float32:
        v.SetFloat(float64(value))
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagFloat as " + kind.String())
    }
    return
}

func readTagDouble(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readFloat64(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.Float32, reflect.Float64:
        v.SetFloat(value)
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
    default:
        return errors.New("cannot parse tagDouble as " + kind.String())
    }
    return
}

func readTagString(reader *bufio.Reader, v reflect.Value) (err error) {
    value, err := readString(reader)
    if err != nil {
        return
    }

    switch kind := v.Kind(); kind {
    case reflect.String:
        v.SetString(value)
        break
    case reflect.Interface:
        v.Set(reflect.ValueOf(value))
        break
    default:
        return errors.New("cannot parse tagString as " + kind.String())
    }
    return
}

func readTagList(reader *bufio.Reader, v reflect.Value) (err error) {
    listType, err := readTagType(reader)
    if err != nil {
        return
    }

    length, err := readInt32(reader)
    if err != nil {
        return
    }

    if length < 0 {
        length = 0
    }

    switch v.Kind() {
    case reflect.Interface:
        v.Set(reflect.ValueOf(make([]interface{}, length)))
        break
    case reflect.Slice:
        v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))
        break
    case reflect.Array:
        if arrayLength := v.Len(); arrayLength < int(length) {
            return fmt.Errorf("size mismatch in tagList: want=%d, available=%d", arrayLength, length)
        }
    }

    for i := 0; i < int(length); i++ {
        err = readValue(reader, listType, v.Index(i))
        if err != nil {
            return
        }
    }
    return
}

func readTagCompoundStruct(reader *bufio.Reader, v reflect.Value) (err error) {
    for {
        var cmpTagType namedTagType
        var cmpTagName string

        cmpTagType, err = readTagType(reader)
        if err != nil {
            return
        }

        if cmpTagType == tagEnd {
            break
        }

        cmpTagName, err = readString(reader)
        if err != nil {
            return
        }

        for i := 0; i < v.NumField(); i++ {
            f := v.Type().Field(i)
            tagName := f.Tag.Get("nbt")
            if tagName == "-" {
                continue
            }

            if tagName == "" {
                tagName = f.Name
            }

            if tagName == cmpTagName {
                err = readValue(reader, cmpTagType, v.Field(i))
                if err != nil {
                    return
                }
                break
            }
        }
    }
    return
}

func readTagCompoundMap(reader *bufio.Reader, v reflect.Value) (err error) {
    if v.Type().Key().Kind() != reflect.String {
        return errors.New("map key should be of type string")
    }

    if v.IsNil() {
        v.Set(reflect.MakeMap(v.Type()))
    }

    for {
        var cmpTagType namedTagType
        var cmpTagName string

        cmpTagType, err = readTagType(reader)
        if err != nil {
            return
        }

        if cmpTagType == tagEnd {
            break
        }

        cmpTagName, err = readString(reader)
        if err != nil {
            return
        }

        if err != nil {
            return err
        }

        var val interface{}

        switch cmpTagType {
        case tagByte:
            val = byte(0)
            break
        case tagShort:
            val = int16(0)
            break
        case tagInt:
            val = int32(0)
            break
        case tagLong:
            val = int64(0)
            break
        case tagFloat:
            val = float32(0)
            break
        case tagDouble:
            val = float64(0)
            break
        case tagString:
            val = ""
            break
        case tagList:
            val = make([]interface{}, 0)
            break
        case tagCompound:
            val = make(map[string]interface{})
            break
        case tagByteArray:
            val = make([]byte, 0)
            break
        case tagIntArray:
            val = make([]int32, 0)
            break
        case tagLongArray:
            val = make([]int64, 0)
            break
        }

        err = readValue(reader, cmpTagType, reflect.ValueOf(&val).Elem())
        if err != nil {
            return err
        }
        v.SetMapIndex(reflect.ValueOf(cmpTagName), reflect.ValueOf(val))
    }
    return
}

func readTagByteArray(reader *bufio.Reader, v reflect.Value) (err error) {
    b, err := readByteSlice(reader)
    if err != nil {
        return
    }

    if v.Type() == reflect.TypeOf(b) {
        v.SetBytes(b)
    } else {
        v.Set(reflect.ValueOf(b))
    }
    return
}

func readTagIntArray(reader *bufio.Reader, v reflect.Value) (err error) {
    length, err := readInt32(reader)
    if err != nil {
        return
    }

    if length < 0 {
        length = 0
    }

    if v.Type().Kind() != reflect.Slice {
        return errors.New("slice required")
    }

    elemType := v.Type().Elem().Kind()
    if elemType != reflect.Int && elemType != reflect.Int32 {
        return errors.New("slice of int or int32 type required")
    }
    v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))
    for i := 0; i < int(length); i++ {
        var val int32
        val, err = readInt32(reader)
        if err != nil {
            return
        }
        v.Index(i).SetInt(int64(val))
    }
    return
}

func readTagLongArray(reader *bufio.Reader, v reflect.Value) (err error) {
    length, err := readInt32(reader)
    if err != nil {
        return
    }

    if length < 0 {
        length = 0
    }

    if v.Type().Kind() != reflect.Slice {
        return errors.New("slice required")
    }

    elemType := v.Type().Elem().Kind()
    if elemType != reflect.Int && elemType != reflect.Int64 {
        return errors.New("slice of int or int64 type required")
    }
    v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))
    for i := 0; i < int(length); i++ {
        var val int64
        val, err = readInt64(reader)
        if err != nil {
            return
        }
        v.Index(i).SetInt(val)
    }
    return
}
