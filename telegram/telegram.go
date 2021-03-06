package telegram

import (
	telegram_bot "github.com/padchin/telegram-bot-api"
	"github.com/padchin/utility"
	"log"
	"strconv"
)

func LoadMessageID(obj *map[string][]int) error {
	err := utility.JSONLoad(obj, "message_id.json")
	if err != nil {
		log.Printf("LoadMessageID error: %v", err)
		return err
	}
	return nil
}

func LoadMessageIDPass(obj *[]int) error {
	err := utility.JSONLoad(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("LoadMessageID error: %v", err)
		*obj = []int{}
		return err
	}
	return nil
}

func DumpMessageIDPass(obj *[]int) error {
	err := utility.JSONDump(obj, "message_id_pass.json")
	if err != nil {
		log.Printf("DumpMessageIDPass error: %v", err)
		return err
	}
	return nil
}

func DumpMessageID(obj *map[string][]int) error {
	err := utility.JSONDump(obj, "message_id.json")
	if err != nil {
		log.Printf("DumpMessageID error: %v", err)
		return err
	}
	return nil
}

// удаление ненужных сообщений
func DeleteWasteMessages(i_user_id int64, bot **telegram_bot.BotAPI, b_is_passphrase bool) {
	s_user_id := strconv.Itoa(int(i_user_id))
	x_message_id_to_delete := make(map[string][]int)

	var ai_message_id_to_delete_pass_phrase []int
	if !b_is_passphrase {
		_ = LoadMessageID(&x_message_id_to_delete)
	} else {
		_ = LoadMessageIDPass(&ai_message_id_to_delete_pass_phrase)
	}
	if !b_is_passphrase {
		if len(x_message_id_to_delete[s_user_id]) > 0 {
			for _, i_message_id := range x_message_id_to_delete[s_user_id] {
				msg_delete := telegram_bot.NewDeleteMessage(i_user_id, i_message_id)
				_, _ = (*bot).Send(msg_delete)
			}
			x_message_id_to_delete[s_user_id] = nil
		}
	} else {
		if len(ai_message_id_to_delete_pass_phrase) > 0 {
			for _, i_message_id := range ai_message_id_to_delete_pass_phrase {
				msg_delete := telegram_bot.NewDeleteMessage(i_user_id, i_message_id)
				_, _ = (*bot).Send(msg_delete)
			}
			ai_message_id_to_delete_pass_phrase = nil
		}
	}
	if !b_is_passphrase {
		_ = DumpMessageID(&x_message_id_to_delete)
	} else {
		_ = DumpMessageIDPass(&ai_message_id_to_delete_pass_phrase)
	}
}

type Keyboard telegram_bot.InlineKeyboardMarkup

func SendMessageWithMarkup(i_user_id int64, s_message string, bot **telegram_bot.BotAPI, markup Keyboard, b_delete_previous bool) {
	if b_delete_previous {
		DeleteWasteMessages(i_user_id, bot, false)
	}
	s_user_id := strconv.Itoa(int(i_user_id))
	msg := telegram_bot.NewMessage(
		i_user_id,
		s_message,
	)
	msg.ReplyMarkup = markup
	msg.ParseMode = "markdown"
	reply, _ := (*bot).Send(msg)
	UpdateMIArrays(s_user_id, reply.MessageID, false)
}

func SendMessage(i_user_id int64, s_message string, bot **telegram_bot.BotAPI, b_delete_previous bool) {
	if b_delete_previous {
		DeleteWasteMessages(i_user_id, bot, false)
	}
	s_user_id := strconv.Itoa(int(i_user_id))
	msg := telegram_bot.NewMessage(
		i_user_id,
		s_message,
	)
	msg.ParseMode = "markdown"
	reply, _ := (*bot).Send(msg)
	UpdateMIArrays(s_user_id, reply.MessageID, false)
}

// UpdateMIArrays обновляет существующие массивы с идентификаторами сообщений для удаления
// b_is_passphrase == true если это сообщения, связанные с паролями
func UpdateMIArrays(s_user_id string, i_message_id int, b_is_passphrase bool) {
	x_message_id_to_delete := make(map[string][]int)
	var ai_message_id_to_delete_pass_phrase []int

	if !b_is_passphrase {
		_ = LoadMessageID(&x_message_id_to_delete)
	} else {
		_ = LoadMessageIDPass(&ai_message_id_to_delete_pass_phrase)
	}
	if !b_is_passphrase {
		x_message_id_to_delete[s_user_id] = append(x_message_id_to_delete[s_user_id], i_message_id)
	} else {
		ai_message_id_to_delete_pass_phrase = append(ai_message_id_to_delete_pass_phrase, i_message_id)
	}
	if !b_is_passphrase {
		_ = DumpMessageID(&x_message_id_to_delete)
	} else {
		_ = DumpMessageIDPass(&ai_message_id_to_delete_pass_phrase)
	}
}
