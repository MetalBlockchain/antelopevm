package service

import (
	"github.com/gin-gonic/gin"
)

var handlers = make(map[string]func(VM) gin.HandlerFunc)

func RegisterHandler(path string, handler func(VM) gin.HandlerFunc) {
	if _, ok := handlers[path]; ok {
		panic("conflicting handler found")
	}

	handlers[path] = handler
}

func GetHandlers() map[string]func(VM) gin.HandlerFunc {
	return handlers
}
