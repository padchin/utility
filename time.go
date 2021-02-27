package utility

import (
	"fmt"
	"time"
)

// LocalTimeFormat возвращает форматированную строку
func LocalTimeFormat(unix_time int) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc
	dt_time := time.Unix(int64(unix_time), 0).Local()
	return fmt.Sprintf(
		"%02d.%02d.%d %02d:%02d:%02d",
		dt_time.Day(),
		dt_time.Month(),
		dt_time.Year(),
		dt_time.Hour(),
		dt_time.Minute(),
		dt_time.Second(),
	)
}
