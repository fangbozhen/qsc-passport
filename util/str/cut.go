package str

import "fmt"

func StrCut(s string, maxlen int) string {
	txt := []rune(s) //需要分割的字符串内容，将它转为字符，然后取长度。

	if len(txt) <= maxlen {
		return s
	}

	return fmt.Sprintf("%s...", string(txt[:maxlen]))
}
