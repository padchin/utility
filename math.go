package utility

import "math"

// Round преобразует число с плавающей точкой в число с плавающей точкой
// с <precision> значащими цифрами после запятой. Получается один из способов округления к ближайшему значащему разряду.
func Round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output+math.Copysign(0.5, num*output))) / output
}
