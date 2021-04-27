package utility

import (
	"bufio"
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

type ReporterOptions struct {
	Bot            *telegram.BotAPI
	Locker         *sync.Mutex
	LogFileName    string
	LastReported   *time.Time
	ReportInterval time.Duration
	ReportMessage  string
}

// Reporter возвращает true, если сообщение об ошибке опубликовано в логах и в Telegram, при наличии связи. Если
// указан интервал 0, то ошибка публикуется в любом случае. Если интервал не 0, то нужно указать ссылку на время
// последней публикации, которое при удачной публикации изменяется на текущее.
func Reporter(options ReporterOptions) bool {
	if options.ReportInterval == 0 || !options.LastReported.IsZero() && time.Since(*options.LastReported) > options.ReportInterval {
		if !options.LastReported.IsZero() {
			*options.LastReported = time.Now()
		}

		if options.Locker != nil {
			options.Locker.Lock()
			defer options.Locker.Unlock()
		}

		logFile, err := os.OpenFile(options.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

		if err != nil {
			log.Print(err)
		}

		log.SetOutput(logFile)
		log.Print(options.ReportMessage)

		err = logFile.Close()

		if err != nil {
			log.Print(err)
		}

		fmt.Println(options.ReportMessage) //nolint:forbidigo

		if options.Bot != nil {
			go tb.SendMessage(iAdminChatID, options.ReportMessage, options.Bot, false)
		}

		return true
	}

	return false
}

// LogFileReduceByTime убирает все данные из лога, которые старше установленного периода от текущей даты.
func LogFileReduceByTime(logFile string, logDuration time.Duration, locker *sync.Mutex) error {
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

	err = origFile.Close()
	if err != nil {
		return fmt.Errorf("error closing original log file: %v", err)
	}

	_, err = exec.Command("mv", logFile, logFile+".bak").Output()

	if err != nil {
		return fmt.Errorf("error creating backup of original log file: %v", err)
	}

	_, err = exec.Command("mv", logFile+".new", logFile).Output()

	if err != nil {
		return fmt.Errorf("error moving new log file to original log file: %v", err)
	}

	return nil
}
