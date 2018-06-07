package resource

import (
	"time"
	"sync"
	"encoding/json"
	"strconv"
	"github.com/andersfylling/disgord/request"
	"errors"
	"github.com/andersfylling/snowflake"
)

const (
	_ int = iota
	MessageActivityTypeJoin
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)
const (
	MessageTypeDefault = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
)

func NewMessage() *Message {
	return &Message{}
}

func NewDeletedMessage() *DeletedMessage {
	return &DeletedMessage{}
}

type DeletedMessage struct {
	ID        snowflake.ID `json:"id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

// https://discordapp.com/developers/docs/resources/channel#message-object-message-application-structure
type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  string       `json:"cover_image"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Name        string       `json:"name"`
}

// Message https://discordapp.com/developers/docs/resources/channel#message-object-message-structure
type Message struct {
	ID              snowflake.ID       `json:"id"`
	ChannelID       snowflake.ID       `json:"channel_id"`
	Author          *User              `json:"author"`
	Content         string             `json:"content"`
	Timestamp       time.Time          `json:"timestamp"`
	EditedTimestamp time.Time          `json:"edited_timestamp"` // ?
	Tts             bool               `json:"tts"`
	MentionEveryone bool               `json:"mention_everyone"`
	Mentions        []*User            `json:"mentions"`
	MentionRoles    []snowflake.ID     `json:"mention_roles"`
	Attachments     []*Attachment      `json:"attachments"`
	Embeds          []*ChannelEmbed    `json:"embeds"`
	Reactions       []*Reaction        `json:"reactions"` // ?
	Nonce           snowflake.ID       `json:"nonce"`     // ?, used for validating a message was sent
	Pinned          bool               `json:"pinned"`
	WebhookID       snowflake.ID       `json:"webhook_id"` // ?
	Type            uint               `json:"type"`
	Activity        MessageActivity    `json:"activity"`
	Application     MessageApplication `json:"application"`

	sync.RWMutex `json:"-"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m.ID.Empty() {
		return []byte("{}"), nil
	}

	//TODO: remove copying of mutex
	return json.Marshal(Message(*m))
}

func (m *Message) Delete() {}
func (m *Message) Update() {}
func (m *Message) Send()   {}

func (m *Message) AddReaction(reaction *Reaction) {}
func (m *Message) RemoveReaction(id snowflake.ID) {}

// ReqGetChannelMessagesParams https://discordapp.com/developers/docs/resources/channel#get-channel-messages-query-string-params
// TODO: ensure limits
type ReqGetChannelMessagesParams struct {
	Around snowflake.ID `urlparam:"around,omitempty"`
	Before snowflake.ID `urlparam:"before,omitempty"`
	After  snowflake.ID `urlparam:"after,omitempty"`
	Limit  int          `urlparam:"limit,omitempty"`
}

// getQueryString this ins't really pretty, but it works.
func (params *ReqGetChannelMessagesParams) getQueryString() string {
	seperator := "?"
	query := ""

	if !params.Around.Empty() {
		query += seperator + params.Around.String()
		seperator = "&"
	}

	if !params.Before.Empty() {
		query += seperator + params.Before.String()
		seperator = "&"
	}

	if !params.After.Empty() {
		query += seperator + params.After.String()
		seperator = "&"
	}

	if params.Limit > 0 {
		query += seperator + strconv.Itoa(params.Limit)
	}

	return query
}

// ReqGetChannelMessages [GET] 	Returns the messages for a channel. If operating on a guild channel, this
// 								endpoint requires the 'VIEW_CHANNEL' permission to be present on the current
// 								user. If the current user is missing the 'READ_MESSAGE_HISTORY' permission
// 								in the channel then this will return no messages (since they cannot read
// 								the message history). Returns an array of message objects on success.
// Endpoint				   		/channels/{channel.id}/messages
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#get-channel-messages
// Reviewed				   		2018-06-07
// Comment				   		The before, after, and around keys are mutually exclusive, only one may
// 								be passed at a time. see ReqGetChannelMessagesParams.
func ReqGetChannelMessages(client request.DiscordGetter, channelID snowflake.ID, params *ReqGetChannelMessagesParams) ([]*Message, error) {
	if channelID.Empty() {
		return nil, errors.New("channelID must be set to get channel messages")
	}
	query := ""
	if params != nil {
		query += params.getQueryString()
	}

	ratelimiter := "/channels/" + channelID.String()
	endpoint := ratelimiter + "/messages" + query
	var messages []*Message
	_, err := client.Get(ratelimiter, endpoint, messages)
	return messages, err
}

// ReqGetChannelMessage [GET] 	Returns a specific message in the channel. If operating on a guild channel,
// 								this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present
// 								on the current user. Returns a message object on success.
// Endpoint				   		/channels/{channel.id}/message/{message.id}
// Rate limiter [MAJOR]	   		/channels/{channel.id}
// Discord documentation   		https://discordapp.com/developers/docs/resources/channel#get-channel-message
// Reviewed				   		2018-06-07
// Comment				   		-
func ReqGetChannelMessage(client request.DiscordGetter, channelID, messageID snowflake.ID) (*Message, error) {
	if channelID.Empty() {
		return nil, errors.New("channelID must be set to get channel messages")
	}
	if messageID.Empty() {
		return nil, errors.New("messageID must be set to get a specific message from a channel")
	}

	ratelimiter := "/channels/" + channelID.String()
	endpoint := ratelimiter + "/message/" + messageID.String()
	var message *Message
	_, err := client.Get(ratelimiter, endpoint, message)
	return message, err
}

type ReqCreateMessageParams struct {
	Content     string        `json:"content"`
	Nonce       snowflake.ID  `json:"nonce,omitempty"`
	Tts         bool          `json:"tts,omitempty"`
	File        interface{}   `json:"file,omitempty"` // TODO: what is this supposed to be?
	Embed       *ChannelEmbed `json:"embed,omitempty"` // embedded rich content
	PayloadJSON string        `json:"payload_json,omitempty"`
}

// ReqCreateChannelMessage [POST]	Post a message to a guild text or DM channel. If operating on a guild channel,
// 									this endpoint requires the 'SEND_MESSAGES' permission to be present on the
// 									current user. If the tts field is set to true, the SEND_TTS_MESSAGES permission
// 									is required for the message to be spoken. Returns a message object. Fires a
// 									Message Create Gateway event. See message formatting for more information on
// 									how to properly format messages.
// 									The maximum request size when sending a message is 8MB.
// Endpoint				   			/channels/{channel.id}/messages
// Rate limiter [MAJOR]	   			/channels/{channel.id}
// Discord documentation   			https://discordapp.com/developers/docs/resources/channel#create-message
// Reviewed				   			2018-06-07
// Comment				   			Before using this endpoint, you must connect to and identify with a gateway
// 									at least once. This endpoint supports both JSON and form data bodies. It does
// 									require multipart/form-data requests instead of the normal JSON request type
// 									when uploading files. Make sure you set your Content-Type to multipart/form-data
// 									if you're doing that. Note that in that case, the embed field cannot be used,
// 									but you can pass an url-encoded JSON body as a form value for payload_json.
// TODO: replace message with a message param struct(!)
func ReqCreateChannelMessage(client request.DiscordPoster, channelID snowflake.ID, params *ReqCreateMessageParams) (*Message, error) {
	if channelID.Empty() {
		return nil, errors.New("channelID must be set to get channel messages")
	}
	if params == nil {
		return nil, errors.New("message must be set")
	}

	ratelimiter := "/channels/" + channelID.String()
	endpoint := ratelimiter + "/messages"
	var generatedMessage *Message
	_, err := client.Post(ratelimiter, endpoint, generatedMessage, params)
	return generatedMessage, err
}
