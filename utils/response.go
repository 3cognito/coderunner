package utils

import "github.com/gin-gonic/gin"

type ServerResponse struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	ErrorMessage *string     `json:"error_message,omitempty"`
	ErrorData    interface{} `json:"error_data,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	response := ServerResponse{
		Success:      true,
		Message:      message,
		ErrorMessage: nil,
		Data:         data,
	}

	ctx.JSON(statusCode, response)
}

func ErrorResponse(ctx *gin.Context, statusCode int, message, errorMessage string, errorData interface{}) {
	response := ServerResponse{
		Success:      false,
		Message:      message,
		ErrorMessage: &errorMessage,
		ErrorData:    errorData,
	}

	ctx.JSON(statusCode, response)
}
