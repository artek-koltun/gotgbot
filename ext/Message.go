package ext

import (
	"encoding/json"
)

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   User   `json:"user"`
}

type Audio struct {
	FileId    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
	MimeType  string `json:"mime_type"`
	FileSize  int    `json:"file_size"`
}

type Document struct {
	FileId   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Video struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type Voice struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type VideoNote struct {
	FileId   string    `json:"file_id"`
	Length   int       `json:"length"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	FileSize int       `json:"file_size"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      int    `json:"user_id"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareId string   `json:"foursquare_id"`
}

type PreCheckoutQuery struct {
	Id               string    `json:"id"`
	From             *User     `json:"from"`
	Currency         string    `json:"currency"`
	TotalAmount      int       `json:"total_amount"`
	InvoicePayload   string    `json:"invoice_payload"`
	ShippingOptionId string    `json:"shipping_option_id"`
	OrderInfo        OrderInfo `json:"order_info"`
}

type Message struct {
	Bot Bot
	MessageId             int                `json:"message_id"`
	From                  *User              `json:"from"`
	Date                  int                `json:"date"`
	Chat                  *Chat              `json:"chat"`
	ForwardFrom           *User              `json:"forward_from"`
	ForwardFromChat       *Chat              `json:"forward_from_chat"`
	ForwardFromMessageId  int                `json:"forward_from_message_id"`
	ForwardSignature      string             `json:"forward_signature"`
	ForwardDate           int                `json:"forward_date"`
	ReplyToMessage        *Message           `json:"reply_to_message"`
	EditDate              int                `json:"edit_date"`
	AuthorSignature       string             `json:"author_signature"`
	Text                  string             `json:"text"`
	Entities              []MessageEntity    `json:"entities"`
	Audio                 *Audio             `json:"audio"`
	Document              *Document          `json:"document"`
	Game                  *Game              `json:"game"`
	Photo                 *PhotoSize         `json:"photo"`
	Sticker               *Sticker           `json:"sticker"`
	Video                 *Video             `json:"video"`
	Voice                 *Voice             `json:"voice"`
	VideoNote             *VideoNote         `json:"video_note"`
	NewChatMembers        []User             `json:"new_chat_members"`
	Caption               string             `json:"caption"`
	Contact               *Contact           `json:"contact"`
	Location              *Location          `json:"location"`
	Venue                 *Venue             `json:"venue"`
	LeftChatMember        *User              `json:"left_chat_member"`
	NewChatTitle          string             `json:"new_chat_title"`
	NewChatPhoto          []PhotoSize        `json:"new_chat_phot"`
	DeleteChatPhoto       bool               `json:"delete_chat_photo"`
	GroupChatCreated      bool               `json:"group_chat_created"`
	SupergroupChatCreated bool               `json:"supergroup_chat_created"`
	ChannelChatCreated    bool               `json:"channel_chat_created"`
	MigrateToChatId       int                `json:"migrate_to_chat_id"`
	MigrateFromChatId     int                `json:"migrate_from_chat_id"`
	PinnedMessage         *Message           `json:"pinned_message"`
	Invoice               *Invoice           `json:"invoice"`
	SuccessfulPayment     *SuccessfulPayment `json:"successful_payment"`
}

func (b Bot) Message(chatId int, text string) Message {
	return Message{Bot: b}
}

func (b Bot) ParseMessage(message json.RawMessage) *Message {
	var mess Message
	json.Unmarshal(message, &mess)
	mess.Bot = b
	return &mess
}

func (m Message) ReplyMessage(text string) (*Message, error) {
	return m.Bot.ReplyMessage(m.Chat.Id, text, m.MessageId)
}

func (m Message) ReplyAudioStr(audio string) (*Message, error) {
	return m.Bot.ReplyAudioStr(m.Chat.Id, audio, m.MessageId)
}

func (m Message) ReplyDocumentStr(document string) (*Message, error) {
	return m.Bot.ReplyDocumentStr(m.Chat.Id, document, m.MessageId)
}

func (m Message) ReplyLocation(latitude float64, longitude float64) (*Message, error) {
	return m.Bot.ReplyLocation(m.Chat.Id, latitude, longitude, m.MessageId)
}

func (m Message) ReplyPhotoStr(photo string) (*Message, error) {
	return m.Bot.ReplyPhotoStr(m.Chat.Id, photo, m.MessageId)
}

func (m Message) ReplyStickerStr(sticker string) (*Message, error) {
	return m.Bot.ReplyStickerStr(m.Chat.Id, sticker, m.MessageId)
}

func (m Message) ReplyVenue(latitude float64, longitude float64, title string, address string) (*Message, error) {
	return m.Bot.ReplyVenue(m.Chat.Id, latitude, longitude, title, address, m.MessageId)
}

func (m Message) ReplyVideoStr(video string) (*Message, error) {
	return m.Bot.ReplyVideoStr(m.Chat.Id, video, m.MessageId)
}

func (m Message) ReplyVideoNoteStr(videoNote string) (*Message, error) {
	return m.Bot.ReplyVideoNoteStr(m.Chat.Id, videoNote, m.MessageId)
}

func (m Message) ReplyVoiceStr(voice string) (*Message, error) {
	return m.Bot.ReplyVoiceStr(m.Chat.Id, voice, m.MessageId)
}

func (m Message) Delete() (bool, error) {
	return m.Bot.DeleteMessage(m.Chat.Id, m.MessageId)
}

func (m Message) Forward(chatId int) (*Message, error) {
	return m.Bot.ForwardMessage(chatId, m.Chat.Id, m.MessageId)
}