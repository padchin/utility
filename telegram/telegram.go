package telegram

import (
	telegram "github.com/padchin/telegram-bot-api"
	"github.com/padchin/utility"
	"log"
	"strconv"
)

func loadMessageID(obj *map[string][]int) error {
	err := utility.JSONLoad(obj, "message_id.json")
	if err != nil {
		log.Printf("loadMessageID error: %v", err)
		return err
	}
	return nil
}

func loadMessageIDPass(obj *[]int) error {
	err := utility.JSONLoad(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("loadMessageID error: %v", err)
		*obj = []int{}
		return err
	}
	return nil
}

func dumpMessageIDPass(obj *[]int) error {
	err := utility.JSONDump(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("dumpMessageIDPass error: %v", err)
		return err
	}
	return nil
}

func dumpMessageID(obj *map[string][]int) error {
	err := utility.JSONDump(obj, "message_id.json")
	if err != nil {
		log.Printf("dumpMessageID error: %v", err)
		return err
	}
	return nil
}

// удаление ненужных сообщений
func DeletePreviousMessages(userID int64, bot **telegram.BotAPI, isPassphrase bool) {
	sUserID := strconv.Itoa(int(userID))
	m := make(map[string][]int)
	var p []int

	if !isPassphrase {
		_ = loadMessageID(&m)
	} else {
		_ = loadMessageIDPass(&p)
	}
	if !isPassphrase {
		if len(m[sUserID]) > 0 {
			for _, iMessageID := range m[sUserID] {
				msgDelete := telegram.NewDeleteMessage(userID, iMessageID)
				_, _ = (*bot).Send(msgDelete)
			}
			m[sUserID] = nil
		}
	} else {
		if len(p) > 0 {
			for _, iMessageID := range p {
				msgDelete := telegram.NewDeleteMessage(userID, iMessageID)
				_, _ = (*bot).Send(msgDelete)
			}
			p = nil
		}
	}
	if !isPassphrase {
		_ = dumpMessageID(&m)
	} else {
		_ = dumpMessageIDPass(&p)
	}
}

type Keyboard telegram.InlineKeyboardMarkup

func EditMessageWithMarkup(iUserID int64, iMessageID int, s_message string, bot **telegram.BotAPI, markup telegram.InlineKeyboardMarkup) {
	msg := telegram.NewEditMessageTextAndMarkup(
		iUserID,
		iMessageID,
		s_message,
		telegram.InlineKeyboardMarkup(markup),
	)
	msg.ParseMode = "markdown"
	_, _ = (*bot).Send(msg)
}

func SendMessageWithMarkup(iUserID int64, message string, bot **telegram.BotAPI, markup telegram.InlineKeyboardMarkup, deletePrevious bool) {
	if deletePrevious {
		DeletePreviousMessages(iUserID, bot, false)
	}
	sUserID := strconv.Itoa(int(iUserID))
	msg := telegram.NewMessage(
		iUserID,
		message,
	)
	msg.ReplyMarkup = markup
	msg.ParseMode = "markdown"
	reply, _ := (*bot).Send(msg)
	storeKeyboardMessageID(iUserID, reply.MessageID)
	UpdateMIArrays(sUserID, reply.MessageID, false)
}

func SendMessage(isUserID int64, message string, bot **telegram.BotAPI, deletePrevious bool) {
	if deletePrevious {
		DeletePreviousMessages(isUserID, bot, false)
	}
	sUserID := strconv.Itoa(int(isUserID))
	msg := telegram.NewMessage(
		isUserID,
		message,
	)
	msg.ParseMode = "markdown"
	reply, _ := (*bot).Send(msg)
	UpdateMIArrays(sUserID, reply.MessageID, false)
}

// UpdateMIArrays обновляет существующие массивы с идентификаторами сообщений для удаления.
// is_passphrase == true если это сообщения, связанные с паролями, которые отправляются только админу.
func UpdateMIArrays(userIDKey string, messageID int, isPassphrase bool) {
	m := make(map[string][]int)
	var p []int

	if !isPassphrase {
		_ = loadMessageID(&m)
	} else {
		_ = loadMessageIDPass(&p)
	}
	if !isPassphrase {
		m[userIDKey] = append(m[userIDKey], messageID)
	} else {
		p = append(p, messageID)
	}
	if !isPassphrase {
		_ = dumpMessageID(&m)
	} else {
		_ = dumpMessageIDPass(&p)
	}
}

func storeKeyboardMessageID(iUserID int64, iMessageID int) {
	m := make(map[int64]int)
	_ = utility.JSONLoad(&m, "message_id_kb.json")
	m[iUserID] = iMessageID
	err := utility.JSONDump(&m, "message_id_kb.json")
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func GetKeyboardMessageID(iUserID int64) int {
	m := make(map[int64]int)
	err := utility.JSONLoad(&m, "message_id_kb.json")
	if err != nil {
		log.Printf("%v", err)
		return 0
	}
	return m[iUserID]
}
