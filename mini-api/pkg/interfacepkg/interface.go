package interfacepkg

import "encoding/json"

// Marshal convert interface json to string
func Marshal(data interface{}) (res string) {
	name, err := json.Marshal(data)
	if err != nil {
		return res
	}
	res = string(name)

	return res
}
