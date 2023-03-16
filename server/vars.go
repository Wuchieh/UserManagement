package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	srv *http.Server
	r   *gin.Engine
)
