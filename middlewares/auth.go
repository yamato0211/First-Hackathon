package middlewares

import (
	"context"
	"first-hackathon/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authString string

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			return
		}

		bearer := "Bearer "
		auth = auth[len(bearer):]

		validate, err := utils.JwtValidate(context.Background(), auth)
		if err != nil || !validate.Valid {
			http.Error(c.Writer, "Invalid token", http.StatusForbidden)
			return
		}

		customClaim, _ := validate.Claims.(*utils.JwtCustomClaim)

		ctx := context.WithValue(c.Request.Context(), authString("auth"), customClaim)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CtxValue(ctx context.Context) *utils.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*utils.JwtCustomClaim)
	return raw
}
