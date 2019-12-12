package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const input = `
["Hello world", 10, false]
`

type Notification struct {
	Message  string
	Priority uint8
	Critical bool
}

func (n *Notification) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Message, &n.Priority, &n.Critical}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}

func main() {
	var n Notification
	if err := json.Unmarshal([]byte(input), &n); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", n)
}
