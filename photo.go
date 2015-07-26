package telegram

import (
	"fmt"

	"bytes"
	"io"

	"mime/multipart"
	"net/http"

	"encoding/json"

	"github.com/AmandaCameron/go-telegram/api"
)

type photoReply struct {
	tgbotapi.PhotoConfig

	bot *Bot
}

type uploadPhotoReply struct {
	*Message
	caption string

	r   io.Reader
	bot *Bot
}

type stickerReply struct {
	tgbotapi.StickerConfig

	bot *Bot
}

// UplaodPhoto uploads a new photo to the service, and sends it as a reply to
// this message.
func (msg *Message) UploadPhoto(r io.Reader, caption string) Uploadable {
	return &uploadPhotoReply{
		Message: msg.ReplyWith(""),

		r:       r,
		caption: caption,

		bot: msg.bot,
	}
}

// PhotoReply sends an already-uploaded photo and sends it as a reply to this
// message.
func (msg *Message) PhotoReply(fileID, caption string) Sendable {
	return photoReply{
		PhotoConfig: tgbotapi.PhotoConfig{
			ChatID:           msg.Chat.ID,
			ReplyToMessageID: msg.replyID,

			UseExistingPhoto: true,

			FileID:  fileID,
			Caption: caption,
		},

		bot: msg.bot,
	}
}

func (pr photoReply) Send() error {
	_, err := pr.bot.api.SendPhoto(pr.PhotoConfig)

	return err
}

// StickerReply replys to this message with a sticker, as defined by fileID.
func (msg *Message) StickerReply(fileID string) Sendable {
	return stickerReply{
		StickerConfig: tgbotapi.StickerConfig{
			ChatID: msg.Chat.ID,

			ReplyToMessageID: msg.MessageID,

			UseExistingSticker: true,

			FileID: fileID,
		},

		bot: msg.bot,
	}
}

func (upl *uploadPhotoReply) Upload() (string, error) {
	if upl.r == nil {
		return "", fmt.Errorf("Invalid reader.")
	}

	bot := upl.bot

	buff := bytes.NewBuffer(nil)
	body := multipart.NewWriter(buff)

	body.WriteField("chat_id", fmt.Sprintf("%d", upl.Chat.ID))
	body.WriteField("caption", upl.caption)
	body.WriteField("reply_to_message_id", fmt.Sprintf("%d", upl.replyID))

	wr, err := body.CreateFormFile("photo", "photo.png")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(wr, upl.r); err != nil {
		return "", err
	}

	body.Close()

	req, err := http.NewRequest("POST",
		"https://api.telegram.org/bot"+bot.api.Token+"/sendPhoto",
		buff,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", body.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var apiResp tgbotapi.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	if !apiResp.Ok {
		return "", fmt.Errorf("Error code %d: %s", apiResp.ErrorCode, apiResp.Description)
	}

	resp.Body.Close()

	var msg tgbotapi.Message

	if err := json.NewDecoder(bytes.NewReader([]byte(apiResp.Result))).Decode(&msg); err != nil {
		return "", err
	}

	var top int32
	fileID := ""

	for _, photo := range msg.Photo {
		if photo.Width*photo.Height > top {
			top = photo.Width * photo.Height

			fileID = photo.FileID
		}
	}

	if fileID == "" {
		return "", fmt.Errorf("Unspecified error happen in the telegram bot API.")
	}

	return fileID, nil
}

func (sr stickerReply) Send() error {
	_, err := sr.bot.api.SendSticker(sr.StickerConfig)

	return err
}
