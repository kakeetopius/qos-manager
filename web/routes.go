package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *ServerCtx) loginPost(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username != "admin" || password != "1234" {
		ctx.Error(ServerError{
			StatusCode: http.StatusOK,
			Err:        errors.New(" Invalid username or password"),
		})

		return
	}

	ctx.Header("HX-Redirect", "/ping")
	ctx.Status(http.StatusOK)
}

func (app *ServerCtx) login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"Title": "Login - QoS Manager",
	})
}
