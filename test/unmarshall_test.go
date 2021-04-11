package test

import (
    "bufio"
    "bytes"
    "github.com/junglemc/nbt"
    "reflect"
    "testing"
)

func TestUnmarshall(t *testing.T) {
    tests := []struct {
        name        string
        tagBytes    []byte
        wantTagName string
        want        interface{}
        wantErr     bool
    }{
        {
            name:        "unnamed root compound tag",
            tagBytes:    UnnamedRootCompoundBytes,
            wantTagName: "",
            want: UnnamedRootCompound{
                ByteTag:   0xFF,
                StringTag: "hello, world",
            },
            wantErr: false,
        },
        {
            name:        "bananrama",
            tagBytes:    BananramaBytes,
            wantTagName: "",
            want:        BananramaStruct,
            wantErr:     false,
        },
        {
            name:        "bigtest",
            tagBytes:    BigTestBytes,
            wantTagName: "Level",
            want: BigTest{
                LongTest:   9223372036854775807,
                ShortTest:  32767,
                StringTest: "HELLO WORLD THIS IS A TEST STRING \xc3\x85\xc3\x84\xc3\x96!",
                FloatTest:  0.49823147058486938,
                IntTest:    2147483647,
                NCT: BigTestNCT{
                    Egg: BigTestNameAndFloat32{
                        Name:  "Eggbert",
                        Value: 0.5,
                    },
                    Ham: BigTestNameAndFloat32{
                        Name:  "Hampus",
                        Value: 0.75,
                    },
                },
                ListTest: []int64{11, 12, 13, 14, 15},
                ListTest2: [2]BigTestCompound{
                    {
                        Name:      "Compound tag #0",
                        CreatedOn: 1264099775885,
                    },
                    {
                        Name:      "Compound tag #1",
                        CreatedOn: 1264099775885,
                    },
                },
                ByteTest:      127,
                ByteArrayTest: BigTestByteArray(),
                DoubleTest:    0.49312871321823148,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            reader := bufio.NewReader(bytes.NewReader(tt.tagBytes))

            actualRaw := reflect.New(reflect.TypeOf(tt.want))

            tagName, err := nbt.Unmarshall(reader, actualRaw.Interface())
            if (err != nil) != tt.wantErr {
                t.Errorf("Unmarshall() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(tt.want, actualRaw.Elem().Interface()) {
                t.Errorf("tags not equal")
                return
            }
        })
    }
}
