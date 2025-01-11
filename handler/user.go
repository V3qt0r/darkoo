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


type UserHandler struct {
	userService models.IUserService
}


func NewUserHandler(UserService models.IUserService) *UserHandler {
	h := &UserHandler{ userService: UserService }
	return h
}


func (h *UserHandler) RegisterUser(c *gin.Context) {
	var request api.RegisterPayload

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error binding data")
		return
	}

	request.Sanitize()
	registerUserPayload := &models.User{
		Email: 		request.Email,
		Password: 	request.Password,
	}

	user, err := h.userService.RegisterUser(registerUserPayload)

	if err != nil {
		log.Print("Error registering user")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Error registering user", nil))
		return
	}

	c.JSON(http.StatusCreated, api.NewResponse(http.StatusCreated, "Successful", user))
}


func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)

	h.getUserById(c, userId)
}


func (h *UserHandler) GetLoggedInUser(c *gin.Context) {
	userDetails, _ := c.Get("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := int(userDetails.(*middleware.User).ID)

	h.getUserById(c, userId)
}


func (h *UserHandler) getUserById(c *gin.Context, id int) {
	user, err := h.userService.GetUserById(id)

	if err != nil {
		log.Print("Error getting user by ID")
		e := apperrors.NewNotFound("Unable to find user with ID %v\n", strconv.Itoa(id))
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Message, gin.H {
			"error": e,
		}))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", user))
}


func (h *UserHandler) GetUserByEmailOrUserName(c *gin.Context) {
	var request map[string] interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Print("Invalid user payload. Email or Username")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Invalid user payload", nil))
		return
	}

	email, ok := request["email"].(string)
	if !ok {
		e := apperrors.NewInternal()
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, e.Message, nil))
		return
	}

	user, err := h.userService.GetUserByEmailOrUserName(email)

	if err != nil {
		log.Print("Unable to get user by email or username")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to get user by email or username.", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", user))
}


func (h *UserHandler) GetUsersByGroupId(c *gin.Context) {
	id := c.Param("id")
	limit := c.Query("limit")
	page := c.Query("page")

	groupId, _ := strconv.Atoi(id)
	limitValue, _ := strconv.Atoi(limit)
	pageValue, _ := strconv.Atoi(page)


	users, err := h.userService.GetUsersByGroupId(groupId, limitValue, pageValue)

	if err != nil {
		log.Print("Unable to get users in group")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Unable to get users in group", nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", users))
}


func (h *UserHandler) UpdateUser(c *gin.Context) {
	userDetails, _ := c.Get("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	var request api.UpdateUserPayload
	if ok := api.BindData(c, &request); !ok {
		e := apperrors.NewBadRequest("Invalid user payload")
		c.JSON(e.Status(), e)
		return
	}

	request.Sanitize()
	user := request.ToEntity()
	user.ID = userId

	err := h.userService.UpdateUser(user)

	if err != nil {
		log.Print("Update user failed!")
		e := apperrors.GetAppError(err, "Update user failed!")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var request api.UpdatePasswordPayload

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error deserializing json data")
		return
	}

	if request.Password != request.ConfirmPassword {
		log.Print("Passwords do not match")
		c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Failed", "Passwords do not match"))
		return
	}

	userDetails, _ := c.Get("id")

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := int(userDetails.(*middleware.User).ID)

	err := h.userService.UpdatePassword(userId, request.Password)

	if err != nil {
		log.Print("Update password failed!")
		e := apperrors.GetAppError(err, "Update password failed!")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) ConfirmPassword(c *gin.Context) {
	userDetails, _ := c.Get("id")
	var request api.ConfirmPasswordPayload

	if userDetails == nil {
		log.Print("User not authenticated")
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, "User not authenticated", nil))
		return
	}

	userId := int(userDetails.(*middleware.User).ID)

	if ok := api.BindData(c, &request); !ok {
		log.Print("Error deserializing json data")
		return
	}

	err := h.userService.ConfirmPassword(userId, request.Password)

	if err != nil {
		log.Print("Password could not be verified")
		e := apperrors.GetAppError(err, "Password could not be verified")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), gin.H{ "isValid" : false}))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) SendOneTimePassword(c *gin.Context) {
	var request api.InitOnetimePasswordPayload
	if ok := api.BindData(c, &request); !ok {
		return
	}

	err := h.userService.InitLoginWithOneTimePassword(request.Email)

	if err != nil {
		log.Print(err)
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) EnrollTOTP(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
				"Error retrieving user details", nil))
		return
	}
	userId := userDetails.(*middleware.User).ID

	totpQRCode, err := h.userService.EnrollTOTP(int(userId))

	if err != nil {
		e := apperrors.GetAppError(err, "User totp enrollment failed")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}
	c.Data(http.StatusOK, "image/png", totpQRCode)
}


func (h *UserHandler) VerifyTOTP(c *gin.Context) {
	var request models.VerifyTOTPRequest
	if ok := api.BindData(c, &request); !ok {
		return
	}

	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
			"Error getting user details", nil))
		return
	}

	userId := userDetails.(*middleware.User).ID

	err := h.userService.VerifyTOTP(int(userId), request)
	if err != nil {
		e := apperrors.GetAppError(err, "Error verifying user totp")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}


func (h *UserHandler) DisableTOTP(c *gin.Context) {
	userDetails, _ := c.Get("id")
	if userDetails == nil {
		c.JSON(http.StatusInternalServerError, api.NewResponse(http.StatusInternalServerError, 
			"Error getting user details", nil))
		return
	}
	userId := userDetails.(*middleware.User).ID

	err := h.userService.DisableTOTP(int(userId))
	if err != nil {
		e := apperrors.GetAppError(err, "Cannot disable totp, please try again")
		c.JSON(e.Status(), api.NewResponse(e.Status(), e.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.NewResponse(http.StatusOK, "Successful", nil))
}