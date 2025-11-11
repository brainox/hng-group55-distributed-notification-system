package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total       int  `json:"total"`
	Limit       int  `json:"limit"`
	Page        int  `json:"page"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

func Success(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func SuccessWithMeta(c *gin.Context, code int, data interface{}, message string, meta *Meta) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}

func Error(c *gin.Context, code int, err error, message string) {
	c.JSON(code, Response{
		Success: false,
		Error:   err.Error(),
		Message: message,
	})
}

func ErrorMessage(c *gin.Context, code int, errorMsg string, message string) {
	c.JSON(code, Response{
		Success: false,
		Error:   errorMsg,
		Message: message,
	})
}

func CalculateMeta(total, limit, page int) *Meta {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &Meta{
		Total:       total,
		Limit:       limit,
		Page:        page,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}
