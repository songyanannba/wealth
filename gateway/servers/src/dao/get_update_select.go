package dao

// MakeUpdateData 组合map数据
func MakeUpdateData(key string, val interface{}) map[string]interface{} {
	var updateMap map[string]interface{}
	updateMap = make(map[string]interface{})
	updateMap[key] = val
	return updateMap
}

func AddUpdateData(key string, val interface{}, m map[string]interface{}) map[string]interface{} {
	m[key] = val
	return m
}

//func MakeUpdateDataS(keys []string, val []interface{}) map[string]interface{} {
//	var updateMap map[string]interface{}
//	updateMap = make(map[string]interface{})
//
//	if len(keys) != len(val) {
//		return updateMap
//	}
//
//	for i, key := range keys {
//		updateMap[key] = val[i]
//	}
//	return updateMap
//}
