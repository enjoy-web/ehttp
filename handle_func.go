package ehttp

import "github.com/gin-gonic/gin"

// HandlerFunc the callback func to handler HTTP Request
//   error -- the result of checking the parameter in HTTP Request (if has no error, error == nil)
type HandlerFunc func(*gin.Context, error)
