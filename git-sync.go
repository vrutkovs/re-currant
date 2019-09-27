package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Env) reload(c *gin.Context) {
	// TODO: restart git-sync sidecar
	c.JSON(http.StatusOK, "{'message': 'ok'}")
}
