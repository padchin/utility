package telegram

import (
	"github.com/padchin/utility"
	"log"
)

func DumpMessageID(obj *interface{}, is_password bool) error {
	if !is_password {
		err := utility.JSONDump(obj, "message_id.json")
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	} else {
		err := utility.JSONDump(obj, "message_id_pass.json")
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}
	return nil
}

func LoadMessageID(obj *interface{}, is_password bool) error {
	if !is_password {
		err := utility.JSONLoad(obj, "message_id.json")
		if err != nil {
			log.Printf("%v", err)
			*obj = make(map[string][]int)
			return err
		}
	} else {
		err := utility.JSONLoad(obj, "message_id_pass.json")
		if err != nil {
			log.Printf("%v", err)
			*obj = []int{}
			return err
		}
	}
	return nil
}
