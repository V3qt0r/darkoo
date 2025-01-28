package main

import (
	"log"
	"os"
	"net/http"

	"darkoo/middleware"
	"darkoo/repository"
	services "darkoo/services"
	"darkoo/websocket"


	dhandlers "darkoo/handler"
	ddb "darkoo/datasources"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	ginEngine := gin.Default()
	ginEngine.MaxMultipartMemory = 8 << 20

	_, isEnv := os.LookupEnv("ENV")
	if !isEnv {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Error loading .env file")
		}
	}

	darkooDB, err := ddb.InitDS()

	if err != nil {
		log.Printf("Error on application startup: %v\n", err)
	}

	userRepository := repository.NewUserRepository(darkooDB.DB)
	groupRepository := repository.NewGroupRepository(darkooDB.DB)
	messageRepository := repository.NewMessageRepository(darkooDB.DB)

	userService := services.NewUserService(userRepository)
	groupService := services.NewGroupService(groupRepository)
	messageService := services.NewMessageService(messageRepository)

	userHandler := dhandlers.NewUserHandler(userService)
	groupHandler := dhandlers.NewGroupHandler(groupService)
	messageHandler := dhandlers.NewMessageHandler(messageService)


	jwtMiddleware, err := middleware.MiddleWare(userService)

	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	errInit := jwtMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error: " + errInit.Error())
	}

	ginEngine.GET("/", func(c *gin.Context) {
		c.File("cors.html")
	})

	ginEngine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Page not found",
		})
	})

	ginEngine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})


	userGroup := ginEngine.Group("/api/user")

	userGroup.POST("/register", userHandler.RegisterUser)
	userGroup.POST("/login", jwtMiddleware.LoginHandler)
	userGroup.GET("/refresh", jwtMiddleware.RefreshHandler)
	userGroup.POST("/one-time-login", jwtMiddleware.LoginHandler)
	userGroup.POST("/send-onetime-password", userHandler.SendOneTimePassword)
	userGroup.POST("/logout", jwtMiddleware.LogoutHandler)
	

	userAuthRoutes := ginEngine.Group("/api/users").Use(jwtMiddleware.MiddlewareFunc())
	userAuthRoutes.GET("/:id", userHandler.GetUserById)
	userAuthRoutes.GET("/email-username", userHandler.GetUserByEmailOrUserName)
	userAuthRoutes.GET("/group/:id", userHandler.GetUsersByGroupId)
	userAuthRoutes.PUT("/self", userHandler.UpdateUser)
	userAuthRoutes.PUT("/update-password", userHandler.UpdatePassword)
	userAuthRoutes.PUT("/confirm-password", userHandler.ConfirmPassword)
	userAuthRoutes.POST("/enroll/totp", userHandler.EnrollTOTP)
	userAuthRoutes.POST("/verify/totp", userHandler.VerifyTOTP)
	userAuthRoutes.POST("/disable/totp", userHandler.DisableTOTP)
	userAuthRoutes.GET("/self", userHandler.GetLoggedInUser)
	userAuthRoutes.PUT("/image-num", userHandler.UpdateUserImageNum)
	userAuthRoutes.PUT("/join-group/:id", userHandler.JoinGroup)


	groupGroup := ginEngine.Group("/api/groups").Use(jwtMiddleware.MiddlewareFunc())
	groupGroup.POST("/create", groupHandler.CreateGroup)
	groupGroup.PUT("/update", groupHandler.UpdateGroup)
	groupGroup.GET("/:id", groupHandler.GetGroupById)
	groupGroup.GET("/self", groupHandler.GetGroupsByUserId)
	groupGroup.DELETE("/:id", groupHandler.DeleteGroupById)
	groupGroup.PUT("/ban/:group_id/users/:user_id", groupHandler.BanUserFromGroup)
	groupGroup.PUT("/unban/:group_id/users/:user_id", groupHandler.UnBanUserFromGroup)

	
	messageGroup := ginEngine.Group("/api/messages").Use(jwtMiddleware.MiddlewareFunc())
	messageGroup.POST("/send/groups/:id", messageHandler.SendMessage)
	messageGroup.GET("/groups/:id", messageHandler.GetMessagesInGroup)
	messageGroup.GET("/groups/user/:id", messageHandler.GetUserMessagesInGroup)
	messageGroup.DELETE("/:message_id/groups/:group_id", messageHandler.DeleteMessage)
	messageGroup.PUT("/:id", messageHandler.UpdateMessage)
	messageGroup.GET("/:id", messageHandler.GetMessageById)

	hub := websocket.NewHub(messageService, userService)
	go hub.Start()

	// Register WebSocket endpoint using gin
	ginEngine.GET("/ws", func(c *gin.Context) {
		// Extract user ID and group ID from the request
		userID := c.MustGet("id").(string) // Assuming you've set it in middleware
		groupID := c.DefaultQuery("groupId", "")

		if userID == "" || groupID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId and groupId are required"})
			return
		}

		// Pass the user ID and group ID to the WebSocket handler
		websocket.HandleWebSocket(hub, c.Writer, c.Request)
	})



	ginEngine.Run(":" + os.Getenv("PORT"))
		
}
