package str

import "strconv"

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
