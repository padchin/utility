package math

import "math"

// Round преобразует число с плавающей точкой в число с плавающей точкой
// с <precision> значащими цифрами после запятой. Получается один из способов округления к ближайшему значащему разряду.
func Round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output+math.Copysign(0.5, num*output))) / output
}

// FastPower возвращает x в степени y. Возвращает 0 для отрицательных степеней.
func FastPower(x, y int) int {
	if y < 0 {
		return 0
	}

	switch {
	case y == 0:
		return 1
	case y == 1:
		return x
	case y%2 == 0:
		return FastPower(x*x, y/2)
	default:
		return x * FastPower(x*x, (y-1)/2)
	}
}
