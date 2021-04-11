package test

import (
    "bufio"
    "bytes"
    "github.com/junglemc/nbt"
    "os"
    "path/filepath"
    "testing"
)

func TestMarshal(t *testing.T) {
    tests := []struct {
        name    string
        tagName string
        tag     interface{}
        want    []byte
        wantErr bool
    }{
        {
            name:    "unnamed root compound tag",
            tagName: "",
            tag: UnnamedRootCompound{
                ByteTag:   0xFF,
                StringTag: "hello, world",
            },
            want:    UnnamedRootCompoundBytes,
            wantErr: false,
        },
        {
            name:    "bananrama",
            tagName: "hello world",
            tag:     BananramaStruct,
            want:    BananramaBytes,
            wantErr: false,
        },
        {
            name:    "bigtest",
            tagName: "Level",
            tag: BigTest{
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
            want:    BigTestBytes,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf := new(bytes.Buffer)
            err := nbt.Marshall(bufio.NewWriter(buf), tt.tagName, tt.tag)
            if (err != nil) != tt.wantErr {
                t.Errorf("Marshall() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            b := buf.Bytes()

            path, err := os.MkdirTemp("", "nbt")
            f, _ := os.Create(filepath.Join(path, "bigtest_go.nbt"))
            w := bufio.NewWriter(f)
            _, _ = w.Write(b)
            _ = w.Flush()

            if !bytes.Equal(b, tt.want) {
                t.Errorf("got:\n[% 2x]\nwant:\n[% 2x]", b, tt.want)
            }
        })
    }
}
