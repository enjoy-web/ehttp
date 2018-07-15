package rest

import "github.com/gin-gonic/gin"

type HandlerFunc func(*gin.Context, error)
