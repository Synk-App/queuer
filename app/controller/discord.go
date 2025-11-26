package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"synk/gateway/app/model"
)

type Discord struct {
	model     *model.Discord
	postModel *model.Posts
}

type HandleDiscordSendResponse struct {
	Resource ResponseHeader                        `json:"resource"`
	Posts    map[int]HandleDiscordSendDataResponse `json:"posts"`
}

type HandleDiscordSendDataResponse struct {
	Resource ResponseHeader                `json:"resource"`
	Post     HandleDiscordSendInfoResponse `json:"post"`
	Raw      any                           `json:"raw"`
}

type HandleDiscordSendInfoResponse struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
	WebhookId string `json:"webhook_id"`
}

type HandleDiscordSendRequest struct {
	Posts []int `json:"posts"`
}

func NewDiscord(db *sql.DB) *Discord {
	discord := Discord{
		model:     model.NewDiscord(db),
		postModel: model.NewPosts(db),
	}

	return &discord
}

func (d *Discord) HandleSend(w http.ResponseWriter, r *http.Request) {
	SetJsonContentType(w)

	response := HandleDiscordSendResponse{
		Resource: ResponseHeader{
			Ok: true,
		},
		Posts: map[int]HandleDiscordSendDataResponse{},
	}

	publisherUrl := strings.TrimSuffix(os.Getenv("PUBLISHER_ENDPOINT"), "/")

	if publisherUrl == "" {
		response.Resource.Ok = false
		response.Resource.Error = "Publisher URL not set"

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusInternalServerError)

		return
	}

	bodyContent, bodyErr := io.ReadAll(r.Body)

	if bodyErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "error on read message body"

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusBadRequest)

		return
	}

	var post HandleDiscordSendRequest

	jsonErr := json.Unmarshal(bodyContent, &post)

	if jsonErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format"

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusBadRequest)

		return
	}

	if len(post.Posts) == 0 {
		response.Resource.Ok = false
		response.Resource.Error = "`posts` can not be empty"

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusBadRequest)

		return
	}

	postsDb, postsDbErr := d.postModel.List(post.Posts)

	if postsDbErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = postsDbErr.Error()

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusInternalServerError)

		return
	}

	for _, postDb := range postsDb {

	}

	payload := map[string]string{"content": post.Message}
	jsonPayload, jsonPayloadErr := json.Marshal(payload)

	if jsonPayloadErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format on sending message"

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusBadRequest)

		return
	}

	respMessage, errMessage := http.Post(post.WebhookUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if errMessage != nil {
		response.Resource.Ok = false
		response.Resource.Error = errMessage.Error()

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, http.StatusBadRequest)

		return
	}

	defer respMessage.Body.Close()

	bodyBytes, _ := io.ReadAll(respMessage.Body)

	response.Raw = string(bodyBytes)

	var responseContent HandleDiscordPublishDataResponse

	json.Unmarshal(bodyBytes, &responseContent)

	if responseContent.Id == "" {
		response.Resource.Ok = false
		response.Resource.Error = responseContent.Message

		WriteErrorResponse(w, response, "/discord/send", response.Resource.Error, respMessage.StatusCode)

		return
	}

	response.Post.ChannelId = responseContent.ChannelId
	response.Post.WebhookId = responseContent.WebhookId
	response.Post.Id = responseContent.Id

	WriteSuccessResponse(w, response)
}
