package internal

import (
	"bufio"
	"errors"
	"math"
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

type Codec struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func (c *Codec) readTagType() (t TagType, err error) {
	tb, err := c.Reader.ReadByte()
	return TagType(tb), err
}

func (c *Codec) readByte() (v byte, err error) {
	return c.Reader.ReadByte()
}

func (c *Codec) readUInt16() (uint16, error) {
	b := make([]byte, 2)
	n, err := c.Reader.Read(b)
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

func (c *Codec) readInt16() (int16, error) {
	uv, err := c.readUInt16()
	if err != nil {
		return 0, err
	}
	return int16(uv), nil
}

func (c *Codec) readUInt32() (uint32, error) {
	b := make([]byte, 4)
	n, err := c.Reader.Read(b)
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

func (c *Codec) readInt32() (int32, error) {
	v, err := c.readUInt32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (c *Codec) readUInt64() (uint64, error) {
	b := make([]byte, 8)
	n, err := c.Reader.Read(b)
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

func (c *Codec) readInt64() (int64, error) {
	v, err := c.readUInt64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func (c *Codec) readFloat32() (float32, error) {
	v, err := c.readUInt32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

func (c *Codec) readFloat64() (float64, error) {
	v, err := c.readUInt64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(v), nil
}

func (c *Codec) readByteSlice() ([]byte, error) {
	length, err := c.readInt32()
	if err != nil {
		return nil, err
	}
	v := make([]byte, length)
	readBytes, err := c.Reader.Read(v)
	if readBytes > int(length) {
		return v, errors.New("read too many bytes")
	}
	if readBytes < int(length) {
		return v, errors.New("read too few bytes")
	}
	return v, nil
}

// TODO: Modified UTF-8 format
func (c *Codec) readString() (string, error) {
	length, err := c.readUInt16()
	if err != nil {
		return "", err
	}
	if length == 0 {
		return "", nil
	}

	v := make([]byte, length)
	readBytes, err := c.Reader.Read(v)
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

func (c *Codec) writeTagType(t TagType) error {
	err := c.writeByte(byte(t))
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeByte(v byte) error {
	err := c.Writer.WriteByte(v)
	if err != nil {
		return err
	}
	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeUInt16(v uint16) error {
	b := []byte{byte(v >> 8), byte(v)}
	_, err := c.Writer.Write(b)
	if err != nil {
		return err
	}

	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeInt16(v int16) error {
	return c.writeUInt16(uint16(v))
}

func (c *Codec) writeUInt32(v uint32) error {
	b := []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	_, err := c.Writer.Write(b)
	if err != nil {
		return err
	}

	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeInt32(v int32) error {
	return c.writeUInt32(uint32(v))
}

func (c *Codec) writeUInt64(v uint64) error {
	b := []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	_, err := c.Writer.Write(b)
	if err != nil {
		return err
	}

	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) writeInt64(v int64) error {
	return c.writeUInt64(uint64(v))
}

func (c *Codec) writeFloat32(v float32) error {
	return c.writeUInt32(math.Float32bits(v))
}

func (c *Codec) writeFloat64(v float64) error {
	return c.writeUInt64(math.Float64bits(v))
}

func (c *Codec) writeByteSlice(v []byte) error {
	err := c.writeInt32(int32(len(v)))
	if err != nil {
		return err
	}
	_, err = c.Writer.Write(v)
	if err != nil {
		return err
	}
	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// TODO: Modified UTF-8 format
func (c *Codec) writeString(v string) error {
	err := c.writeUInt16(uint16(len(v)))
	if err != nil {
		return err
	}
	_, err = c.Writer.Write([]byte(v))
	if err != nil {
		return err
	}
	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
