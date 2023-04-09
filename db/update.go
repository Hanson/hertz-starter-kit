package db

// stupid gorm only accept map[string]interface{}
func ToUpdateMap(m map[string]string) (res map[string]interface{}) {
	res = make(map[string]interface{})
	for k, v := range m {
		res[k] = v
	}
	return
}
