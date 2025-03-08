package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"githuv.com/wmfadel/go_events/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "UnAuthorized",
		})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("invalid token %v", err),
		})
		return
	}

	context.Set("userId", userId)

	context.Next()
}

// func AuthorizeForEvents(context *gin.Context) {
// 	token := context.Request.Header.Get("Authorization")

// 	if token == "" {
// 		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"message": "UnAuthorized",
// 		})
// 		return
// 	}

// 	userId, err := utils.VerifyToken(token)
// 	if err != nil {
// 		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"message": fmt.Sprintf("invalid token %v", err),
// 		})
// 		return
// 	}

// 	eventId := context.Request.URL.Query().Get("id")
// 	db.DB.GetEventById(strconv.ParseInt(eventId, 10, 64))

// 	context.Next()
// }
