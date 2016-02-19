package bencode

import (
    "fmt"
    "unicode/utf8"
)

func EncodeInt(x int) string {
    return fmt.Sprintf("i%de", x)
}

func EncodeList(x ...string) string {
    // TODO(ian): DRY.
    tmp := "l"
    for i := range x { 
        tmp = fmt.Sprintf("%s%s", tmp, x[i])
    }
    tmp = fmt.Sprintf("%se", tmp)
    return tmp
}

func EncodeDictionary(x ...string) string {
    // TODO(ian): DRY.
    tmp := "d"
    for i := range x { 
        tmp = fmt.Sprintf("%s%s", tmp, x[i])
    }
    tmp = fmt.Sprintf("%se", tmp)
    return tmp
}

func EncodeByteString(x string) string {
    return fmt.Sprintf("%d:%s", utf8.RuneCountInString(x), x)
}
