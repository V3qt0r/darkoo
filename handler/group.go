package handler

import (
	"net/http"
	"log"
	"strconv"

	"darkoo/api"
	"darkoo/models"
	"darkoo/apperrors"
	"darkoo/middleware"

	"github.com/gin-gonic/gin"
)


type GroupHandler struct {
	groupService models.IGroupService
}


func NewGroupHandler(GroupService models.IGroupService) *GroupHandler {
	h := &GroupHandler{ groupService : GroupService }
	return h
}


func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var request api.CreateGroupPayload

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error deserializing json data from group handler")
		return
	}

	request.Sanitize()
	createGroupPayload := &models.Group{
		Name: request.Name,
		Description: request.Description,
	}

	group, err := h.groupService.CreateGroup(createGroupPayload)

	if err != nil {
		log.Print("Error creating group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Error creating group", nil))
		return
	}

	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", group))
}


func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	var request api.UpdateGroupPayload
	id := c.Param("id")
	groupId, _ := strconv.Atoi(id)

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error deserializing json data from group handler")
		return
	}

	request.Sanitize()
	group := request.ToEntity()
	group.ID = uint(groupId)

	err := h.groupService.UpdateGroup(group)

	if err != nil {
		log.Print("Update group failed!")
		e := apperrors.GetAppError(err, "Update group failed!")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *GroupHandler) GetGroupById(c *gin.Context) {
	id := c.Param("id")
	groupId, _ := strconv.Atoi(id)

	group, err := h.groupService.GetGroupById(groupId)

	if err != nil {
		log.Print("Error retrieving group by provided ID")
		e := apperrors.NewNotFound("Unable to find group with ID %v\n", strconv.Itoa(groupId))
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H {
			"error": e,
		}))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", group))
}


func (h *GroupHandler) GetGroupsByUserId(c *gin.Context) {
	userDetails, _ := c.Get("id")
	limit := c.Query("limit")
	page := c.Query("page")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := int(userDetails.(*middleware.User).ID)
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)

	groups, err := h.groupService.GetGroupsByUserId(userId, limitValue, pageValue)

	if err != nil {
		log.Print("Unable to get groups this user belongs to")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to get groups this user belongs to", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", groups))
}


func (h *GroupHandler) DeleteGroupById(c *gin.Context) {
	id := c.Param("id")
	groupId, _ := strconv.Atoi(id)
	
	err := h.groupService.DeleteGroupById(groupId)

	if err != nil {
		log.Print("Failed to delete group")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "Failed to delete group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *GroupHandler) BanUserFromGroup(c *gin.Context) {
	gid := c.Param("group_id")
	uid := c.Param("user_id")

	groupId, _ := strconv.Atoi(gid)
	userId, _ := strconv.Atoi(uid)

	err := h.groupService.BanUserFromGroup(groupId, userId)

	if err != nil {
		log.Print("Unable to ban user from this group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to ban user from this group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *GroupHandler) UnBanUserFromGroup(c *gin.Context) {
	gid := c.Param("group_id")
	uid := c.Param("user_id")

	groupId, _ := strconv.Atoi(gid)
	userId, _ := strconv.Atoi(uid)

	err := h.groupService.UnBanUserFromGroup(groupId, userId)

	if err != nil {
		log.Print("Unable to ban user from this group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to ban user from this group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}