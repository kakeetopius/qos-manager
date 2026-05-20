package web

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		var serverError ServerError
		if errors.As(err, &serverError) {
			ctx.JSON(serverError.StatusCode, gin.H{
				"success": false,
				"message": serverError.Error(),
			})
		}
	}
}

func ErrorHandlerHTML() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		var serverError ServerError
		if errors.As(err, &serverError) {
			ctx.HTML(serverError.StatusCode, "fail", gin.H{
				"Error": serverError.Error(),
			})
		}
	}
}
