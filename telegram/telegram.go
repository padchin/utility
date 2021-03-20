package telegram

import (
	telegram "github.com/padchin/telegram-bot-api"
	"github.com/padchin/utility/io"
	"log"
	"strconv"
)

type TButton struct {
	Text         string
	CallbackData string
}

type Keyboard telegram.InlineKeyboardMarkup

// AddButtonsRow метод для создания одного ряда, состоящего из одной или нескольких кнопок.
func (k *Keyboard) AddButtonsRow(buttons ...TButton) {
	var row []telegram.InlineKeyboardButton
	for _, b := range buttons {
		row = append(row, telegram.NewInlineKeyboardButtonData(
			b.Text,
			b.CallbackData,
		))
	}
	(*k).InlineKeyboard = append((*k).InlineKeyboard, row)
}

func loadMessageID(obj *map[string][]int) error {
	err := io.JSONLoad(obj, "message_id.json")
	if err != nil {
		log.Printf("loadMessageID error: %v", err)
		return err
	}
	return nil
}

func loadMessageIDPass(obj *[]int) error {
	err := io.JSONLoad(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("loadMessageID error: %v", err)
		*obj = []int{}
		return err
	}
	return nil
}

func dumpMessageIDPass(obj *[]int) error {
	err := io.JSONDump(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("dumpMessageIDPass error: %v", err)
		return err
	}
	return nil
}

func dumpMessageID(obj *map[string][]int) error {
	err := io.JSONDump(obj, "message_id.json")
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


func EditMessageWithMarkup(iUserID int64, iMessageID int, sMessage string, bot **telegram.BotAPI, markup telegram.InlineKeyboardMarkup, deletePrevious ...bool) {
	if len(deletePrevious) > 0 {
		if deletePrevious[0] {
			DeletePreviousMessages(iUserID, bot, false)
		}
	}
	msg := telegram.NewEditMessageTextAndMarkup(
		iUserID,
		iMessageID,
		sMessage,
		markup,
	)
	msg.ParseMode = "markdown"
	_, _ = (*bot).Send(msg)
}

// SendMessageWithMarkup отправляет сообщение с клавиатурой с использованием экземпляра бота и добавляет messageID в базу для
// последующего удаления. Если требуется удалить предыдущие сообщения, имеющиеся в базе, установить флаг deletePrevious
// в true.
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

// SendMessage отправляет сообщение без клавиатуры с использованием экземпляра бота и добавляет messageID в базу для
// последующего удаления. Если требуется удалить предыдущие сообщения, имеющиеся в базе, установить флаг deletePrevious
// в true. Если нужно исключить сообщение из базы для последующего удаления установить флаг persist в true.
func SendMessage(iUserID int64, sMessage string, bot **telegram.BotAPI, deletePrevious bool, persist ...bool) {
	if deletePrevious {
		DeletePreviousMessages(iUserID, bot, false)
	}
	sUserID := strconv.Itoa(int(iUserID))
	msg := telegram.NewMessage(
		iUserID,
		sMessage,
	)
	msg.ParseMode = "markdown"
	reply, _ := (*bot).Send(msg)
	if len(persist) > 0 {
		if persist[0] {
			return
		}
	}
	UpdateMIArrays(sUserID, reply.MessageID, false)
}

// UpdateMIArrays обновляет существующие массивы с идентификаторами сообщений для удаления.
// isPassphrase - true если это сообщения, связанные с паролями, которые отправляются только админу.
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
	_ = io.JSONLoad(&m, "message_id_kb.json")
	m[iUserID] = iMessageID
	err := io.JSONDump(&m, "message_id_kb.json")
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

// GetKeyboardMessageID возвращает messageID последней выведенной клавиатуры для редактирования.
func GetKeyboardMessageID(iUserID int64) int {
	m := make(map[int64]int)
	err := io.JSONLoad(&m, "message_id_kb.json")
	if err != nil {
		log.Printf("%v", err)
		return 0
	}
	return m[iUserID]
}
