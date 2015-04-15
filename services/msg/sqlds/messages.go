package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
)

type MessageService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const MESSAGES_TABLE = "messages"
const MESSAGE_RECEIPTS_TABLE = "message_receipts"
const MESSAGE_FRAGMENTS_TABLE = "message_fragments"

func NewMessageService(db *sql.DB, sg *ServiceGroup) *MessageService {
	svc := MessageService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *MessageService) InitDB() {
	svc.SG.IDService.CreateDomain(&CreateDomainRequest{nil, "messageids", 1, 2})
	CreateTable(svc.DB, MESSAGES_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"ChannelId bigint NOT NULL REFERENCES channels (Id) ON DELETE CASCADE",
			"SenderId bigint DEFAULT(0) REFERENCES users (Id) ON DELETE SET DEFAULT",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"MsgType TEXT DEFAULT ('')",
			"MsgData TEXT DEFAULT ('')",
			"Status INT DEFAULT (0)",
		})
	CreateTable(svc.DB, MESSAGE_FRAGMENTS_TABLE,
		[]string{
			"MessageId bigint NOT NULL REFERENCES messages (Id) ON DELETE CASCADE",
			"Index INT NOT NULL",
			"Size INT DEFAULT(0)",
			"MimeType TEXT DEFAULT ('')",
			"Body BYTEA DEFAULT('')",
		},
		", CONSTRAINT unique_message_fragment UNIQUE (MessageId, Index)")
	CreateTable(svc.DB, MESSAGE_RECEIPTS_TABLE,
		[]string{
			"MessageId bigint NOT NULL REFERENCES messages (Id) ON DELETE CASCADE",
			"ReceiverId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"ReceivedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Error TEXT DEFAULT('')",
			"Status INT DEFAULT(0)",
		},
		", CONSTRAINT unique_message_receipt UNIQUE (MessageId, ReceiverId)")
}

/**
 * Send messages to particular recipients in this channel.
 * If Channel, message or user is nil an error is returned.
 */
func (svc *MessageService) SaveMessage(request *SaveMessageRequest) error {
	if request.Message.Id == 0 {
		id, err := svc.SG.IDService.NextID(&NextIDRequest{nil, "messageids"})
		if err != nil {
			return err
		}
		err = InsertRow(svc.DB, MESSAGES_TABLE,
			"Id", id,
			"ChannelId", request.Message.Channel.Id,
			"SenderId", request.Message.Sender.Id,
			"Status", request.Message.Status,
			"MsgType", request.Message.MsgType,
			"MsgData", request.Message.MsgData)
		if err == nil {
			request.Message.Id = id
		}
		return err
	} else {
		return UpdateRows(svc.DB, MESSAGES_TABLE,
			fmt.Sprintf("Id = %d", request.Message.Id),
			"ChannelId", request.Message.Channel.Id,
			"SenderId", request.Message.Sender.Id,
			"Status", request.Message.Status,
			"MsgType", request.Message.MsgType,
			"MsgData", request.Message.MsgData)
	}
}

/**
 * Get the messages in a channel for a particular user (optional)
 * If channel does not exist then an empty list is returned.
 * If user is nil then all messages in a channel is returned.
 * Pagination is possible with offset and count.
 * Messages are also ordered by the created time stamp.
 */
func (svc *MessageService) GetMessages(request *GetMessagesRequest) ([]*Message, error) {
	reverse := false
	descending := false
	if request.Count <= 0 {
		request.Count = 50
	}
	if request.Offset < 0 {
		request.Offset = (-request.Offset) - 1
		descending = true
		reverse = true
	}
	whereClause := fmt.Sprintf("WHERE ChannelId = %d", request.Channel.Id)
	if request.Sender != nil {
		whereClause += fmt.Sprintf(" AND SenderId = %d", request.Sender.Id)
	}
	orderClause := "ORDER BY Created"
	if descending {
		orderClause += " DESC"
	}

	query := fmt.Sprintf("SELECT Id, Created, Status, MsgType, MsgData, SenderId from %s %s %s LIMIT %d OFFSET %d", MESSAGES_TABLE, whereClause, orderClause, request.Count, request.Offset)
	rows, err := svc.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	out := make([]*Message, 0, request.Count)
	for rows.Next() {
		var senderId int64 = 0
		msg := &Message{}
		err := rows.Scan(&msg.Id, &msg.Created, &msg.Status, &msg.MsgType, &msg.MsgData, &senderId)
		if err != nil {
			return nil, err
		}
		if request.Sender == nil {
			msg.Sender, err = svc.SG.UserService.GetUser(NewUserById(senderId))
		} else {
			msg.Sender = request.Sender
		}
		out = append(out, msg)
	}
	if reverse {
		start := 0
		end := len(out) - 1
		for start < end {
			tmp := out[start]
			out[start] = out[end]
			out[end] = tmp
			start++
			end--
		}
	}
	return out, nil
}

/**
 * Gets a message by ID
 */
func (svc *MessageService) GetMessageById(request *GetMessageRequest) (*GetMessageResult, error) {
	return nil, nil
}

/**
 * Gets the fragments of a message.
 */
func (svc *MessageService) GetMessageFragments(request *GetMessageRequest) (*GetMessageResult, error) {
	return nil, nil
}

/**
 * Get receipts of a particular message.
 */
func (svc *MessageService) GetMessageReceipts(request *GetMessageRequest) (*GetMessageResult, error) {
	return nil, nil
}

/**
 * Remove a particular message.
 */
func (svc *MessageService) DeleteMessage(request *DeleteMessageRequest) error {
	return DeleteById(svc.DB, MESSAGES_TABLE, request.Message.Id)
}

/**
 * Removes all entries.
 */
func (svc *MessageService) RemoveAllMessages(request *Request) {
	ClearTable(svc.DB, MESSAGES_TABLE)
}
