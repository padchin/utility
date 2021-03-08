package utility

// MaxValueElementsCount возвращает количество одинаковых элементов в массиве, имеющих максимальное значение.
func MaxValueElementsCount(array *[]int) int {
	aiBias := make(map[int]int)
	for _, x := range *array {
		if _, ok := aiBias[x]; ok {
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
	return valuesCount
}
