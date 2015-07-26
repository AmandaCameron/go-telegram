package tgbotapi

import (
	"encoding/json"
)

// APIResponse is a response from the Telegram API with the result stored raw.
type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int32             `json:"error_code"`
	Description string          `json:"description"`
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID int32     `json:"update_id"`
	Message  Message `json:"message"`
}

// User is a user, contained in Message and returned by GetSelf.
type User struct {
	ID        int32    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

// GroupChat is a group chat, and not currently in use.
type GroupChat struct {
	ID    int32    `json:"id"`
	Title string `json:"title"`
}

// UserOrGroupChat is returned in Message, because it's not clear which it is.
type UserOrGroupChat struct {
	ID        int32    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Title     string `json:"title"`
}

// Message is returned by almost every request, and contains data about almost anything.
type Message struct {
	MessageID           int32             `json:"message_id"`
	From                User            `json:"from"`
	Date                int32             `json:"date"`
	Chat                UserOrGroupChat `json:"chat"`
	ForwardFrom         User            `json:"forward_from"`
	ForwardDate         int32             `json:"forward_date"`
	ReplyToMessage      *Message        `json:"reply_to_message"`
	Text                string          `json:"text"`
	Audio               Audio           `json:"audio"`
	Document            Document        `json:"document"`
	Photo               []PhotoSize     `json:"photo"`
	Sticker             Sticker         `json:"sticker"`
	Video               Video           `json:"video"`
	Contact             Contact         `json:"contact"`
	Location            Location        `json:"location"`
	NewChatParticipant  User            `json:"new_chat_participant"`
	LeftChatParticipant User            `json:"left_chat_participant"`
	NewChatTitle        string          `json:"new_chat_title"`
	NewChatPhoto        string          `json:"new_chat_photo"`
	DeleteChatPhoto     bool            `json:"delete_chat_photo"`
	GroupChatCreated    bool            `json:"group_chat_created"`
}

// PhotoSize contains information about photos, including ID and Width and Height.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int32    `json:"width"`
	Height   int32    `json:"height"`
	FileSize int32    `json:"file_size"`
}

// Audio contains information about audio, including ID and Duration.
type Audio struct {
	FileID   string `json:"file_id"`
	Duration int32    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int32    `json:"file_size"`
}

// Document contains information about a document, including ID and a Thumbnail.
type Document struct {
	FileID    string    `json:"file_id"`
	Thumbnail PhotoSize `json:"thumb"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	FileSize  int32       `json:"file_size"`
}

// Sticker contains information about a sticker, including ID and Thumbnail.
type Sticker struct {
	FileID    string    `json:"file_id"`
	Width     int32       `json:"width"`
	Height    int32       `json:"height"`
	Thumbnail PhotoSize `json:"thumb"`
	FileSize  int32       `json:"file_size"`
}

// Video contains information about a video, including ID and duration and Thumbnail.
type Video struct {
	FileID    string    `json:"file_id"`
	Width     int32       `json:"width"`
	Height    int32       `json:"height"`
	Duration  int32       `json:"duration"`
	Thumbnail PhotoSize `json:"thumb"`
	MimeType  string    `json:"mime_type"`
	FileSize  int32       `json:"file_size"`
	Caption   string    `json:"caption"`
}

// Contact contains information about a contact, such as PhoneNumber and UserId.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      string `json:"user_id"`
}

// Location contains information about a place, such as Longitude and Latitude.
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// UserProfilePhotos contains information a set of user profile photos.
type UserProfilePhotos struct {
	TotalCount int32         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

// ReplyKeyboardMarkup allows the Bot to set a custom keyboard.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]string `json:"keyboard"`
	ResizeKeyboard  bool       `json:"resize_keyboard"`
	OneTimeKeyboard bool       `json:"one_time_keyboard"`
	Selective       bool       `json:"selective"`
}

// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

// ForceReply allows the Bot to have users directly reply to it without additional interaction.
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"force_reply"`
}
