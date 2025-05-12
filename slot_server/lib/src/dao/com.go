package dao

// MakeUpdateData 组合map数据
func MakeUpdateData(key string, val interface{}) map[string]interface{} {
	var updateMap map[string]interface{}
	updateMap = make(map[string]interface{})
	updateMap[key] = val
	return updateMap
}
