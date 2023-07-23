package model

import (
	"encoding/base64"
	"fmt"
)

type MessageConnection struct {
	Messages []*Message
	From     int
	To       int
}

func (u *MessageConnection) TotalCount() int {
	return len(u.Messages)
}

func (u *MessageConnection) PageInfo() PageInfo {
	return PageInfo{
		StartCursor: EncodeCursor(u.From),
		EndCursor:   EncodeCursor(u.To - 1),
		HasNextPage: u.To < len(u.Messages),
	}
}

func EncodeCursor(i int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1)))
}
