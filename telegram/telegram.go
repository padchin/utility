package telegram

import (
	"fmt"
	"strconv"

	telegram "github.com/padchin/telegram-bot-api"
	"github.com/padchin/utility/file_operations"
)

const (
	MARKDOWN             = "markdown"
	MESSAGE_ID_JSON      = "message_id.json"
	MESSAGE_ID_KB_JSON   = "message_id_kb.json"
	MESSAGE_ID_PASS_JSON = "message_id_pass.json"
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
	err := file_operations.JSONLoad(obj, MESSAGE_ID_JSON)

	if err != nil {
		return fmt.Errorf("loadMessageID error: %v", err)
	}

	return nil
}

func loadMessageIDPass(obj *[]int) error {
	err := file_operations.JSONLoad(obj, MESSAGE_ID_PASS_JSON)

	if err != nil {
		*obj = []int{}

		return fmt.Errorf("loadMessageID error: %v", err)
	}

	return nil
}

func dumpMessageIDPass(obj *[]int) error {
	err := file_operations.JSONDump(obj, MESSAGE_ID_PASS_JSON)

	if err != nil {
		return fmt.Errorf("dumpMessageIDPass error: %v", err)
	}

	return nil
}

func dumpMessageID(obj *map[string][]int) error {
	err := file_operations.JSONDump(obj, MESSAGE_ID_JSON)

	if err != nil {
		return fmt.Errorf("dumpMessageID error: %v", err)
	}

	return nil
}

// DeletePreviousMessages удаление ненужных сообщений
func DeletePreviousMessages(userID int64, bot *telegram.BotAPI, isPassphrase bool) {
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
			// исключаем активную клавиатуру из списка на удаление
			exMessage := GetKeyboardMessageID(userID)

			for _, iMessageID := range m[sUserID] {
				if iMessageID != exMessage {
					msgDelete := telegram.NewDeleteMessage(userID, iMessageID)
					_, _ = bot.Send(msgDelete)
				}
			}

			m[sUserID] = nil
		}
	} else {
		if len(p) > 0 {
			for _, iMessageID := range p {
				msgDelete := telegram.NewDeleteMessage(userID, iMessageID)
				_, _ = bot.Send(msgDelete)
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

func EditMessageWithMarkup(iUserID int64, iMessageID int, sMessage string, bot *telegram.BotAPI, markup telegram.InlineKeyboardMarkup, deletePrevious ...bool) error {
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

	msg.ParseMode = MARKDOWN
	_, err := bot.Send(msg)

	return err
}

// SendMessageWithMarkup отправляет сообщение с клавиатурой с использованием экземпляра бота и добавляет messageID в базу для
// последующего удаления. Если требуется удалить предыдущие сообщения, имеющиеся в базе, установить флаг deletePrevious
// в true.
func SendMessageWithMarkup(iUserID int64, message string, bot *telegram.BotAPI, markup telegram.InlineKeyboardMarkup, deletePrevious bool) error {
	if deletePrevious {
		DeletePreviousMessages(iUserID, bot, false)
	}

	sUserID := strconv.Itoa(int(iUserID))
	msg := telegram.NewMessage(
		iUserID,
		message,
	)
	msg.ReplyMarkup = markup
	msg.ParseMode = MARKDOWN
	reply, err := (*bot).Send(msg)
	storeKeyboardMessageID(iUserID, reply.MessageID)
	UpdateMIArrays(sUserID, reply.MessageID, false)

	return err
}

// SendMessage отправляет сообщение без клавиатуры с использованием экземпляра бота и добавляет messageID в базу для
// последующего удаления. Если требуется удалить предыдущие сообщения, имеющиеся в базе, установить флаг deletePrevious
// в true. Если нужно исключить сообщение из базы для последующего удаления установить флаг persist в true.
func SendMessage(iUserID int64, sMessage string, bot *telegram.BotAPI, deletePrevious bool, persist ...bool) error {
	if deletePrevious {
		DeletePreviousMessages(iUserID, bot, false)
	}

	sUserID := strconv.Itoa(int(iUserID))
	msg := telegram.NewMessage(
		iUserID,
		sMessage,
	)

	msg.ParseMode = MARKDOWN
	reply, err := bot.Send(msg)

	if len(persist) > 0 {
		if persist[0] {
			return err
		}
	}

	UpdateMIArrays(sUserID, reply.MessageID, false)

	return err
}

func SendDocument(iUserID int64, documentPath string, bot *telegram.BotAPI, deletePrevious bool, persist ...bool) error {
	if deletePrevious {
		DeletePreviousMessages(iUserID, bot, false)
	}

	sUserID := strconv.Itoa(int(iUserID))
	msg := telegram.NewDocumentUpload(
		iUserID,
		documentPath,
	)

	reply, err := bot.Send(msg)

	if len(persist) > 0 {
		if persist[0] {
			return err
		}
	}

	UpdateMIArrays(sUserID, reply.MessageID, false)

	return err
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
	_ = file_operations.JSONLoad(&m, MESSAGE_ID_KB_JSON)
	m[iUserID] = iMessageID
	err := file_operations.JSONDump(&m, MESSAGE_ID_KB_JSON)

	if err != nil {
		return
	}
}

// GetKeyboardMessageID возвращает messageID последней выведенной клавиатуры для редактирования.
func GetKeyboardMessageID(iUserID int64) int {
	m := make(map[int64]int)

	err := file_operations.JSONLoad(&m, MESSAGE_ID_KB_JSON)

	if err != nil {
		return 0
	}

	return m[iUserID]
}
