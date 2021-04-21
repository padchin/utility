package utility

import (
	"fmt"
	telegram "github.com/padchin/telegram-bot-api"
	tb "github.com/padchin/utility/telegram"
	"log"
	"time"
)

const iAdminChatID int64 = 726713220

// ErrorReport возвращает true, если сообщение об ошибке опубликовано в логах и в Telegram, при наличии связи. Если
//указан интервал 0, то ошибка публикуется в любом случае. Если интервал не 0, то нужно указать ссылку на время
//последней публикации, которое при удачной публикации изменяется на текущее.
func ErrorReport(bot *telegram.BotAPI, err string, interval time.Duration, lastReported ...*time.Time) bool {
	if interval == 0 || (len(lastReported) > 0 && time.Since(*lastReported[0]) > interval) {
		if len(lastReported) > 0 {
			*lastReported[0] = time.Now()
		}
		log.Print(err)
		fmt.Println(err)
		if bot != nil {
			go tb.SendMessage(iAdminChatID, err, bot, false)
		}
		return true
	}
	return false
}
