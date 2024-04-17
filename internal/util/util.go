package util

func IsDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func ByteToInt(ch byte) int {
    return int(ch - '0')
}
