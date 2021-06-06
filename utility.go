package utility

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	telegram "github.com/padchin/telegram-bot-api"
	tb "github.com/padchin/utility/telegram"
)

const iAdminChatID int64 = 726713220

// ErrNotPublished возвращается, если сообщение не опубликовано.
var ErrNotPublished = errors.New("сообщение не опубликовано")
var ErrTimeNotSpecified = errors.New("не указана ссылка на время последней публикации")

type ReporterOptions struct {
	// ChatIDs содержит массив из идентификаторов пользователей Telegram, которым будет отправлено уведомление.
	ChatIDs *[]int64
	Bot     *telegram.BotAPI
	// Locker для предотвращения одновременной записи в лог.
	Locker *sync.Mutex
	// Имя файла лога включая полный путь при необходимости.
	LogFileName string
	// Указатель на время последней публикации.
	LastReported *time.Time
	// Интервал между публикациями.
	Interval time.Duration
	// Текст сообщения для публикации.
	Message string
}

// Reporter публикует сообщение в логах и в Telegram (список пользователей указывается в параметрах), при
// наличии экземпляра бота. Если указан интервал 0, то сообщение публикуется в любом случае. Если интервал не 0, то
// нужно указать ссылку на время последней публикации, которое при удачной в публикации в логах изменяется на текущее. В
// этом случае сообщение публикуется не чаще, чем через интервал, указанный в параметрах.
func Reporter(r ReporterOptions) error {
	if r.Interval != 0 && r.LastReported == nil {
		return ErrTimeNotSpecified
	}

	if r.Interval == 0 || r.LastReported != nil && time.Since(*r.LastReported) > r.Interval {
		if r.Locker != nil {
			r.Locker.Lock()
			defer r.Locker.Unlock()
		}

		if len(r.LogFileName) > 0 {
			logFile, err := os.OpenFile(r.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

			if err == nil {
				log.SetOutput(logFile)
				log.Print(r.Message)

				_ = logFile.Close()
			}
		}

		fmt.Println(r.Message) //nolint:forbidigo

		if r.Bot != nil {
			if r.ChatIDs == nil || len(*r.ChatIDs) == 0 {
				// если не указан список пользователей, отправляется только админу
				tb.SendMessage(iAdminChatID, r.Message, r.Bot, false, true)
			} else {
				for _, chat := range *r.ChatIDs {
					tb.SendMessage(chat, r.Message, r.Bot, false, true)
				}
			}
		}

		if r.LastReported != nil {
			*r.LastReported = time.Now()
		}

		return nil
	}

	return ErrNotPublished
}

// LogFileReduceByTime убирает все данные из лога, которые старше установленного периода от текущей даты.
func LogFileReduceByTime(logFile string, logDuration time.Duration, locker *sync.Mutex) (err error) {
	if locker != nil {
		locker.Lock()
		defer locker.Unlock()
	}

	origFile, err := os.Open(logFile)

	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}

	newFile, err := os.Create(logFile + ".new")

	if err != nil {
		return fmt.Errorf("error creating new temporary log file: %v", err)
	}

	writer := bufio.NewWriter(newFile)

	scanner := bufio.NewScanner(origFile)

	for scanner.Scan() {
		// date format 2021/04/24 20:39:03
		line := scanner.Text()

		if len(line) < 20 {
			continue
		}

		dtOfLine, err := time.Parse("2006/01/02 15:04:05", line[:19])

		if err != nil {
			continue
		}

		// если дата строки помещается в установленный период
		if time.Since(dtOfLine) < logDuration {
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("error writung to log file: %v", err)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	err = writer.Flush()

	if err != nil {
		return fmt.Errorf("error flushing new temporary log file: %v", err)
	}

	err = newFile.Close()

	if err != nil {
		return fmt.Errorf("error closing new temporary log file: %v", err)
	}

	_ = origFile.Close()

	err = exec.Command("mv", logFile, logFile+".bak").Run()

	if err != nil {
		return fmt.Errorf("error creating backup of original log file: %v", err)
	}

	err = exec.Command("mv", logFile+".new", logFile).Run()

	if err != nil {
		return fmt.Errorf("error moving new log file to original log file: %v", err)
	}

	return nil
}
