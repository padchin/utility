package utility

// MaxNumsCount вычисляет количество элементов в массиве с максимальным значением
func MaxNumsCount(array *[]int) int {
	ai_bias := make(map[int]int)
	for _, x := range *array {
		if _, ok := ai_bias[x]; ok {
			ai_bias[x]++
		} else {
			ai_bias[x] = 1
		}
	}
	//find maximum
	max_val := 0
	value := 0
	for key, count := range ai_bias {
		if count > max_val {
			max_val = count
			value = key
		}

	}
	return value
}
