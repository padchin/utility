package utility

import (
	"time"
)

// LocalTimeFormatFromUnix возвращает форматированную строку в формате 02.01.2006 15:04:05 без долей секунды
func LocalTimeFormatFromUnix(unix_time int) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc
	dt_time := time.Unix(int64(unix_time), 0).Local()
	return dt_time.Format("02.01.2006 15:04:05")
}
