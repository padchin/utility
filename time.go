package utility

import (
	"time"
)

// LocalTimeFormatFromUnix возвращает форматированную строку в формате 02.01.2006 15:04:05 без долей секунды,
// если не указан формат, или в соответствии с указанным форматом в местном времени (+3)
func LocalTimeFormatFromUnix(unixTime int, format ...string) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc
	t := time.Unix(int64(unixTime), 0).Local()
	if len(format) > 0 {
		return t.Format(format[0])
	} else {
		return t.Format("02.01.2006 15:04:05")
	}
}

// LocalTimeFormatFromUnixNano возвращает форматированную строку в формате 02.01.2006 15:04:05 без долей секунды,
// если не указан формат, или в соответствии с указанным форматом в местном времени (+3)
func LocalTimeFormatFromUnixNano(unixTime int, format ...string) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc
	t := time.Unix(int64(unixTime/1e9), int64(unixTime%1e9)).Local()
	if len(format) > 0 {
		return t.Format(format[0])
	} else {
		return t.Format("02.01.2006 15:04:05")
	}
}
