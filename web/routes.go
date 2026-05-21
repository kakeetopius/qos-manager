package web

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kakeetopius/qosm/internal/core/pam"
)

func (app *ServerCtx) loginPost(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	app.logger.Info("login attempt", "username", username)

	if err := pam.AuthenticateUser(username, password); err != nil {
		app.logger.Error("auth_failed", "username", username, "err", err)
		ctx.Error(ServerError{
			StatusCode: http.StatusOK,
			Err:        fmt.Errorf(" Invalid username or password"),
		})

		return
	}

	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge:   300,  // expires in 5minutes
		HttpOnly: true, // Prevent JavaScript access
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	session.Set("username", username)
	session.Set("role", "administrator")
	session.Save()
	ctx.Header("HX-Redirect", "/dashboard")
	ctx.Status(http.StatusOK)
}

func (app *ServerCtx) login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"Title": "Login - QoS Manager",
	})
}

func (app *ServerCtx) dashboard(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "dashboard", gin.H{
		"Heading":     "Dashboard",
		"Description": "Overview of network traffic and QoS policies",
		"User":        session.Get("username"),
		"Role":        session.Get("role"),
	})
}

func (app *ServerCtx) rules(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "rules", gin.H{
		"Heading":     "Rules",
		"Description": "Define how network traffic should be prioritized or limited",
		"User":        session.Get("username"),
		"Role":        session.Get("role"),
	})
}

func (app *ServerCtx) analytics(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "analytics", gin.H{
		"Heading":     "Analytics",
		"Description": "Network usage insights and QoS effectiveness",
		"User":        session.Get("username"),
		"Role":        session.Get("role"),
	})
}

func (app *ServerCtx) logs(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "logs", gin.H{
		"Heading":     "Logs",
		"Description": "Real-time QoS engine and network activity logs",
		"User":        session.Get("username"),
		"Role":        session.Get("role"),
	})
}

func (app *ServerCtx) settings(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "settings", gin.H{
		"Heading":     "Settings",
		"Description": "Configure QoS engine behavior and system preferences",
		"User":        session.Get("username"),
		"Role":        session.Get("role"),
	})
}

func (app *ServerCtx) logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
