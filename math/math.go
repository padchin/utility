package math

import "math"

// Round преобразует число с плавающей точкой в число с плавающей точкой
// с <precision> значащими цифрами после запятой. Получается один из способов округления к ближайшему значащему разряду.
func Round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output+math.Copysign(0.5, num*output))) / output
}

// FastPower возвращает x в степени n. Возвращает 0 для отрицательных степеней.
func FastPower(x, n int) int {
	switch {
	case n < 0:
		return 0
	case n == 0:
		return 1
	case n == 1:
		return x
	case n%2 == 0:
		return FastPower(x*x, n/2)
	default:
		return x * FastPower(x*x, (n-1)/2)
	}
}
