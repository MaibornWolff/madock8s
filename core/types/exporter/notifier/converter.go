package notifier

func ConvertNotifierMapToStringMap(notifierMap *Map) map[string]string {
	resultMap := make(map[string]string, len(notifierMap.Items))
	for _, item := range notifierMap.Items {
		resultMap[item.Key] = item.Value
	}
	return resultMap
}

func ConvertStringMapToNotifierMap(stringMap map[string]string) *Map {
	resultMap := &Map{}
	for key, value := range stringMap {
		item := &Map_Pair{Key: key, Value: value}
		resultMap.Items = append(resultMap.Items, item)
	}
	return resultMap
}
