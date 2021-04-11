package nbt

import (
    "bufio"
    "errors"
    "github.com/junglemc/nbt/internal"
    "reflect"
)

func Marshall(writer *bufio.Writer, tagName string, v interface{}) error {
    val := reflect.ValueOf(v)
    tagType := internal.TypeOf(val.Type())
    return internal.Marshall(writer, tagName, tagType, val)
}

func Unmarshall(reader *bufio.Reader, v interface{}) (tagName string, err error) {
    val := reflect.ValueOf(v)
    if val.Kind() != reflect.Ptr {
        return "", errors.New("ptr required")
    }

    if tagName, err = internal.Unmarshall(reader, &val); err != nil {
        return "", err
    }
    return tagName, nil
}
