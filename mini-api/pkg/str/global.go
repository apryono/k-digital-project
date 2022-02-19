package str

import (
	"regexp"
	"strconv"
)

// StringToInt ...
func StringToInt(data string) int {
	res, err := strconv.Atoi(data)
	if err != nil {
		res = 0
	}

	return res
}

// ShowString ...
func ShowString(isShow bool, data string) string {
	if isShow {
		return data
	}

	return ""
}

// CheckEmail ...
func CheckEmail(text string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(text)
}
