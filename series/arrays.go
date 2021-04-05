package series

// MaxValueElementsCount возвращает максимальный по значению элемент и их количество.
func MaxValueElementsCount(array *[]int) (int, int) {
	aiBias := make(map[int]int)
	for _, x := range *array {
		if //goland:noinspection GoLinterLocal
		_, ok := aiBias[x]; ok {
			aiBias[x]++
		} else {
			aiBias[x] = 1
		}
	}
	//find maximum
	maxVal := 0
	valuesCount := 0
	for keyValue, count := range aiBias {
		if count > maxVal {
			maxVal = count
			valuesCount = keyValue
		}
	}
	return maxVal, valuesCount
}
