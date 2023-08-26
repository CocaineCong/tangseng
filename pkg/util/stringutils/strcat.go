package stringutils

func StrConcat(strs []string) string {
	length := len(strs)
	buf := make([]byte, 0, length)
	for j := 0; j < length; j++ {
		buf = append(buf, strs[j]...)
	}

	return string(buf)
}
