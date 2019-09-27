package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Env) apply(c *gin.Context) {
	// TODO: run kustomize apply
	c.JSON(http.StatusOK, "{'message': 'ok'}")
}
