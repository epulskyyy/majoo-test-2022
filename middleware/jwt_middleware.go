package middleware

import (
	"github.com/epulskyyy/majoo-test-2022/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTMiddleware(authUseCase usecase.IAuthUseCase)  gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata, err := authUseCase.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "You need to be authorized to access this route")
			c.Abort()
		}
		err = authUseCase.FetchAuth(metadata)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "You need to be authorized to access this route")
			c.Abort()
		}
		c.Next()
	}
}
