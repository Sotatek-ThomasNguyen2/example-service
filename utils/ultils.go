package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func BindQuery(in interface{}, ctx *gin.Context) error {
	err := decoder.Decode(in, ctx.Request.URL.Query())
	return err
}
