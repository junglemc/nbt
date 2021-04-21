package nbt

import (
    "bufio"
    "errors"
    "math"
    "reflect"
)

type TagType byte

const (
    TagEnd TagType = iota
    TagByte
    TagShort
    TagInt
    TagLong
    TagFloat
    TagDouble
    TagByteArray
    TagString
    TagList
    TagCompound
    TagIntArray
    TagLongArray
    TagNone = 0xFF
)

func readTagType(reader *bufio.Reader) (t TagType, err error) {
    tb, err := readByte(reader)
    return TagType(tb), err
}

func readByte(reader *bufio.Reader) (v byte, err error) {
    return reader.ReadByte()
}

func readUInt16(reader *bufio.Reader) (uint16, error) {
    b := make([]byte, 2)
    n, err := reader.Read(b)
    if err != nil {
        return 0, err
    }
    if n > 2 {
        return 0, errors.New("too much data")
    }
    if n < 2 {
        return 0, errors.New("not enough data")
    }
    return uint16(b[0])<<8 | uint16(b[1]), nil
}

func readInt16(reader *bufio.Reader) (int16, error) {
    uv, err := readUInt16(reader)
    if err != nil {
        return 0, err
    }
    return int16(uv), nil
}

func readUInt32(reader *bufio.Reader) (uint32, error) {
    b := make([]byte, 4)
    n, err := reader.Read(b)
    if err != nil {
        return 0, err
    }
    if n > 4 {
        return 0, errors.New("too much data")
    }
    if n < 4 {
        return 0, errors.New("not enough data")
    }
    return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3]), nil
}

func readInt32(reader *bufio.Reader) (int32, error) {
    v, err := readUInt32(reader)
    if err != nil {
        return 0, err
    }
    return int32(v), nil
}

func readUInt64(reader *bufio.Reader) (uint64, error) {
    b := make([]byte, 8)
    n, err := reader.Read(b)
    if err != nil {
        return 0, err
    }
    if n > 8 {
        return 0, errors.New("too much data")
    }
    if n < 8 {
        return 0, errors.New("not enough data")
    }
    return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7]), nil
}

func readInt64(reader *bufio.Reader) (int64, error) {
    v, err := readUInt64(reader)
    if err != nil {
        return 0, err
    }
    return int64(v), nil
}

func readFloat32(reader *bufio.Reader) (float32, error) {
    v, err := readUInt32(reader)
    if err != nil {
        return 0, err
    }
    return math.Float32frombits(v), nil
}

func readFloat64(reader *bufio.Reader) (float64, error) {
    v, err := readUInt64(reader)
    if err != nil {
        return 0, err
    }
    return math.Float64frombits(v), nil
}

func readByteSlice(reader *bufio.Reader) ([]byte, error) {
    length, err := readInt32(reader)
    if err != nil {
        return nil, err
    }
    v := make([]byte, length)
    readBytes, err := reader.Read(v)
    if readBytes > int(length) {
        return v, errors.New("read too many bytes")
    }
    if readBytes < int(length) {
        return v, errors.New("read too few bytes")
    }
    return v, nil
}

// TODO: Modified UTF-8 format
func readString(reader *bufio.Reader) (string, error) {
    length, err := readUInt16(reader)
    if err != nil {
        return "", err
    }
    if length == 0 {
        return "", nil
    }

    v := make([]byte, length)
    readBytes, err := reader.Read(v)
    if err != nil {
        return "", err
    }
    if readBytes > int(length) {
        return "", errors.New("read too many bytes")
    }
    if readBytes < int(length) {
        return "", errors.New("read too few bytes")
    }
    return string(v), nil
}

func writeTagType(writer *bufio.Writer, t TagType) error {
    err := writeByte(writer, byte(t))
    if err != nil {
        return err
    }
    return nil
}

func writeByte(writer *bufio.Writer, v byte) error {
    err := writer.WriteByte(v)
    if err != nil {
        return err
    }
    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeUInt16(writer *bufio.Writer, v uint16) error {
    b := []byte{byte(v >> 8), byte(v)}
    _, err := writer.Write(b)
    if err != nil {
        return err
    }

    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeInt16(writer *bufio.Writer, v int16) error {
    return writeUInt16(writer, uint16(v))
}

func writeUInt32(writer *bufio.Writer, v uint32) error {
    b := []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
    _, err := writer.Write(b)
    if err != nil {
        return err
    }

    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeInt32(writer *bufio.Writer, v int32) error {
    return writeUInt32(writer, uint32(v))
}

func writeUInt64(writer *bufio.Writer, v uint64) error {
    b := []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
    _, err := writer.Write(b)
    if err != nil {
        return err
    }

    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeInt64(writer *bufio.Writer, v int64) error {
    return writeUInt64(writer, uint64(v))
}

func writeFloat32(writer *bufio.Writer, v float32) error {
    return writeUInt32(writer, math.Float32bits(v))
}

func writeFloat64(writer *bufio.Writer, v float64) error {
    return writeUInt64(writer, math.Float64bits(v))
}

func writeByteSlice(writer *bufio.Writer, v []byte) error {
    err := writeInt32(writer, int32(len(v)))
    if err != nil {
        return err
    }
    _, err = writer.Write(v)
    if err != nil {
        return err
    }
    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeInt32Slice(writer *bufio.Writer, v []int32) (err error) {
    n := len(v)

    err = writeInt32(writer, int32(n))
    if err != nil {
        return
    }

    for i := 0; i < n; i++ {
        err = writeInt32(writer, v[i])
        if err != nil {
            return
        }
    }
    return
}

func writeInt64Slice(writer *bufio.Writer, v []int64) (err error) {
    n := len(v)

    err = writeInt32(writer, int32(n))
    if err != nil {
        return
    }

    for i := 0; i < n; i++ {
        err = writeInt64(writer, v[i])
        if err != nil {
            return
        }
    }
    return
}

// TODO: Modified UTF-8 format
func writeString(writer *bufio.Writer, v string) error {
    err := writeUInt16(writer, uint16(len(v)))
    if err != nil {
        return err
    }
    _, err = writer.Write([]byte(v))
    if err != nil {
        return err
    }
    err = writer.Flush()
    if err != nil {
        return err
    }
    return nil
}

func writeList(writer *bufio.Writer, v reflect.Value) (err error) {
    nestedTagType := TypeOf(v.Type().Elem())

    n := v.Len()
    if n <= 0 {
        nestedTagType = TagEnd // Mimic notchian behavior
    }

    err = writeTagType(writer, nestedTagType)
    if err != nil {
        return
    }

    if err = writeInt32(writer, int32(n)); err != nil {
        return
    }

    for i := 0; i < n; i++ {
        err = writeValue(writer, nestedTagType, v.Index(i).Interface())
        if err != nil {
            return
        }
    }
    return
}

func writeCompound(writer *bufio.Writer, v reflect.Value) (err error) {
    if v.Type().Kind() == reflect.Map {
        mapRange := v.MapRange()
        for mapRange.Next() {
            nestedTagType := TypeOf(reflect.TypeOf(mapRange.Value().Interface()))

            err = writeTagType(writer, nestedTagType)
            if err != nil {
                return
            }

            err = writeString(writer, mapRange.Key().String())
            if err != nil {
                return
            }

            err = writeValue(writer, nestedTagType, mapRange.Value().Interface())
            if err != nil {
                return
            }
        }
    } else {
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

            err = writeTagType(writer, nestedTagType)
            if err != nil {
                return
            }

            err = writeString(writer, nestedTagName)
            if err != nil {
                return
            }

            err = writeValue(writer, nestedTagType, v.Field(i).Interface())
            if err != nil {
                return
            }
        }
    }
    return writeTagType(writer, TagEnd)
}
