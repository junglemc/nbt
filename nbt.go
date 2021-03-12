package nbt

import (
	"bufio"
	"errors"
	"github.com/junglemc/nbt/internal"
	"reflect"
)

func Marshal(writer *bufio.Writer, tagName string, v interface{}) error {
	c := new(internal.Codec)
	c.Writer = writer
	val := reflect.ValueOf(v)
	tagType := c.TypeOf(val.Type())
	return c.Marshall(tagName, tagType, val)
}

func Unmarshall(reader *bufio.Reader, v interface{}) (tagName string, err error) {
	c := new(internal.Codec)
	c.Reader = reader
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return "", errors.New("ptr required")
	}

	if tagName, err = c.Unmarshall(&val); err != nil {
		return "", err
	}
	return tagName, nil
}
