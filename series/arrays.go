package series

// MaxValueElementsCount возвращает максимальный по значению элемент и их количество.
func MaxValueElementsCount(array *[]int) (int, int) {
	// aiBias содержит значения массива в качестве ключа и кол-во повторений элемента массива в качестве значения
	aiBias := make(map[int]int)

	for _, x := range *array {
		aiBias[x]++
	}

	// находим максимум
	// берем любой элемент из словаря как максимальный
	val := aiBias[(*array)[0]]
	count := 0

	for v, c := range aiBias {
		if c > val {
			val = c
			count = v
		}
	}

	return val, count
}
