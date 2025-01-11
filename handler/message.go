package handler

import (
	"net/http"
	"log"
	"strconv"

	"darkoo/api"
	"darkoo/middleware"
	"darkoo/models"
	"darkoo/apperrors"

	"github.com/gin-gonic/gin"
)


type MessageHandler struct {
	messageService models.IMessageService
}


func NewMessageHandler(MessageService models.IMessageService) *MessageHandler {
	h := &MessageHandler{ messageService: MessageService }
	return h
}


func (h *MessageHandler) SendMessage(c *gin.Context) {
	var request api.SendMessagePayload
	userDetails, _ := c.Get("id")
	id := c.Param("id")

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error deserializing json data from user handler")
		return
	}

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID
	groupId, _ := strconv.Atoi(id)

	request.Sanitize()
	request.Validate()

	sendMessagePayload := &models.Message{
		Content: request.Content,
		ContentType: request.ContentType,
		AttachmentUrl: request.AttachmentUrl,
		GroupId: uint(groupId),
		UserId: userId,
	}

	message, err := h.messageService.SendMessage(sendMessagePayload)

	if err != nil {
		log.Print("Error sending message")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Error sending message", nil))
		return
	}

	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", message))
}


func (h *MessageHandler) GetMessagesInGroup(c *gin.Context) {
	id := c.Param("id")
	limit := c.Query("limit")
	page := c.Query("page")

	groupId, _ := strconv.Atoi(id)
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	messages, err := h.messageService.GetMessagesInGroup(groupId, limitValue, pageValue)

	if err != nil {
		log.Print("Unable to get messages in this group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to get messages in this group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", messages))
}


func (h *MessageHandler) GetUserMessagesInGroup(c *gin.Context) {
	userDetails, _ := c.Get("id")
	id := c.Param("id")
	limit := c.Query("limit")
	page := c.Query("page")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := int(userDetails.(*middleware.User).ID)
	groupId, _ := strconv.Atoi(id)
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	messages, err := h.messageService.GetUserMessagesInGroup(userId, groupId, limitValue, pageValue)

	if err != nil {
		log.Print("Unable to get user messages in this group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to get user messages in this group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", messages))
}


func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userDetails, _ := c.Get("id")
	mId := c.Param("message_id")
	gId := c.Param("group_id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	messageId, _ := strconv.Atoi(mId)
	groupId, _ := strconv.Atoi(gId)
	userId := int(userDetails.(*middleware.User).ID)

	err := h.messageService.DeleteMessage(messageId, userId, groupId)

	if err != nil {
		log.Print("Failed to delete message")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Failed to delete message", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	var request api.UpdateMessagePayload
	id := c.Param("id")
	userDetails, _ := c.Get("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	request.Validate()

	if ok := api.BindData(c, &request); !ok {
		e := apperrors.NewBadRequest("Invalid message payload")
		c.JSON(e.Status(), e)
		return
	}

	messageId, _ := strconv.Atoi(id)
	userId := userDetails.(*middleware.User).ID



	message := request.ToEntity()
	message.ID = uint(messageId)

	foundMessage, _ := h.messageService.GetMessageById(messageId)

	if foundMessage.UserId != userId {
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "You can only edit your own message", nil))
		return
	}

	err := h.messageService.UpdateMessage(message)

	if err != nil {
		log.Print("Update message failed!")
		e := apperrors.GetAppError(err, "Update message failed!")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *MessageHandler) GetMessageById(c *gin.Context) {
	id := c.Param("id")
	messageId, _ := strconv.Atoi(id)

	message, err := h.messageService.GetMessageById(messageId)

	if err != nil {
		log.Print("Error getting message by ID")
		e := apperrors.NewNotFound("Unable to find message with ID %v\n", id)
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H {
			"error": e,
		}))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", message))
}