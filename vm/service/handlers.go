package service

import (
	"github.com/gin-gonic/gin"
)

var handlers = make(map[string]Handler)

type Handler struct {
	HandlerFunc func(VM) gin.HandlerFunc
	Methods     []string
}

func RegisterHandler(path string, handler Handler) {
	if _, ok := handlers[path]; ok {
		panic("conflicting handler found")
	}

	handlers[path] = handler
}

func GetHandlers() map[string]Handler {
	return handlers
}
