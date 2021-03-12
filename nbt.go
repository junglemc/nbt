package nbt

import (
	"bufio"
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
