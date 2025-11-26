package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"synk/gateway/app/model"
)

type Telegram struct {
	model *model.Telegram
}

type HandleTelegramPublishResponse struct {
	Resource ResponseHeader                    `json:"resource"`
	Post     HandleTelegramPublishInfoResponse `json:"post"`
	Raw      any                               `json:"raw"`
}

type HandleTelegramPublishInfoResponse struct {
	MessageId string `json:"message_id"`
}

type HandleTelegramPublishDataResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Result      struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"result"`
}

type HandleTelegramPublishRequest struct {
	Message  string `json:"message"`
	BotToken string `json:"bot_token"`
	ChatId   string `json:"chat_id"`
}

func NewTelegram(db *sql.DB) *Telegram {
	telegram := Telegram{
		model: model.NewTelegram(db),
	}

	return &telegram
}

func (d *Telegram) HandleSend(w http.ResponseWriter, r *http.Request) {
	SetJsonContentType(w)

	response := HandleTelegramPublishResponse{
		Resource: ResponseHeader{
			Ok: true,
		},
		Post: HandleTelegramPublishInfoResponse{},
	}

	bodyContent, bodyErr := io.ReadAll(r.Body)

	if bodyErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "error on read message body"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	var post HandleTelegramPublishRequest

	jsonErr := json.Unmarshal(bodyContent, &post)

	if jsonErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	post.BotToken = strings.TrimSpace(post.BotToken)
	post.ChatId = strings.TrimSpace(post.ChatId)
	post.Message = strings.TrimSpace(post.Message)

	if post.BotToken == "" || post.ChatId == "" || post.Message == "" {
		response.Resource.Ok = false
		response.Resource.Error = "field `bot_token`, `chat_id` and `message` can not be empty"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	endpointUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", post.BotToken)

	payload := map[string]string{"chat_id": post.ChatId, "text": post.Message}
	jsonPayload, jsonPayloadErr := json.Marshal(payload)

	if jsonPayloadErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format on sending message"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	respMessage, errMessage := http.Post(endpointUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if errMessage != nil {
		response.Resource.Ok = false
		response.Resource.Error = errMessage.Error()

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	defer respMessage.Body.Close()

	bodyBytes, _ := io.ReadAll(respMessage.Body)

	response.Raw = string(bodyBytes)

	var responseContent HandleTelegramPublishDataResponse

	json.Unmarshal(bodyBytes, &responseContent)

	if !responseContent.Ok {
		response.Resource.Ok = false
		response.Resource.Error = responseContent.Description

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, respMessage.StatusCode)

		return
	}

	response.Post.MessageId = strconv.Itoa(responseContent.Result.MessageID)

	WriteSuccessResponse(w, response)
}
