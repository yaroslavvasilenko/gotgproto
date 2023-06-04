package types

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/jaskaur18/gotgproto/errors"
	"github.com/jaskaur18/gotgproto/functions"
)

type Message struct {
	*tg.Message
	ReplyToMessage *Message
	Text           string
	IsService      bool
	Action         tg.MessageActionClass
}

func ConstructMessage(m tg.MessageClass) *Message {
	switch msg := m.(type) {
	case *tg.Message:
		return constructMessageFromMessage(msg)
	case *tg.MessageService:
		return constructMessageFromMessageService(msg)
	case *tg.MessageEmpty:
		return constructMessageFromMessageEmpty(msg)
	}
	return &Message{}
}

func constructMessageFromMessage(m *tg.Message) *Message {
	return &Message{
		Message: m,
		Text:    m.Message,
	}
}

func constructMessageFromMessageEmpty(m *tg.MessageEmpty) *Message {
	return &Message{
		Message: &tg.Message{
			ID:     m.ID,
			PeerID: m.PeerID,
		},
	}
}

func constructMessageFromMessageService(m *tg.MessageService) *Message {
	return &Message{
		Message: &tg.Message{
			Out:         m.Out,
			Mentioned:   m.Mentioned,
			MediaUnread: m.MediaUnread,
			Silent:      m.Silent,
			Post:        m.Post,
			Legacy:      m.Legacy,
			ID:          m.ID,
			Date:        m.Date,
			FromID:      m.FromID,
			PeerID:      m.PeerID,
			ReplyTo:     m.ReplyTo,
			TTLPeriod:   m.TTLPeriod,
		},
		IsService: true,
		Action:    m.Action,
	}
}

func (m *Message) SetRepliedToMessage(ctx context.Context, raw *tg.Client) error {
	replyTo := m.ReplyTo.ReplyToMsgID
	if replyTo == 0 {
		return errors.ErrMessageNotExist
	}
	chatId := functions.GetChatIdFromPeer(m.PeerID)
	msgs, err := functions.GetMessages(ctx, raw, chatId, []tg.InputMessageClass{
		&tg.InputMessageID{
			ID: replyTo,
		},
	})
	if err != nil {
		return err
	}
	m.ReplyToMessage = ConstructMessage(msgs[0])
	return nil
}
