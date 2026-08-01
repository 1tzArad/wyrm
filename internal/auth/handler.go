package auth

import (
	"net/http"

	"github.com/1tzArad/wyrm/pkg/response"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body!")
			return
		}

		if err := service.Register(c.Request.Context(), req.Username, req.Password); err != nil {
			if err == ErrUserExists {
				response.Fail(c, http.StatusConflict, "USERNAME_EXISTS", "username already taken")
				return
			}
			response.InternalFail(c)
			return
		}

		response.OK(c, gin.H{"message": "user created successfully"})
	}
}

func LoginHandler(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body!")
			return
		}

		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials!")
			return
		}

		c.SetCookie("token", token, int(tokenExpiration), "/", "", true, true)
		response.OK(c, gin.H{"message": "user logged in successfully!"})
	}
}
