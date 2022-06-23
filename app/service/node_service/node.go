package node_service

import (
	"github.com/gin-gonic/gin"
	"go-dage-web/app/dao"
	"net/http"
	"path/filepath"
	"strings"
)

func GetList(ctx *gin.Context) {
	list := dao.List{}
	if err := list.Get(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "获取数据失败，请联系管理员",
		})
	} else {
		for i, _ := range list {
			list[i].Date = strings.TrimSuffix(list[i].Date, filepath.Ext(list[i].Date))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Code":  http.StatusOK,
			"Msg":   nil,
			"Items": list,
		})
	}
}

func AddNode(ctx *gin.Context) {
	name := ctx.Query("name")
	item := dao.ListItem{Name: name}
	if succeed, err := item.Add(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code":   http.StatusInternalServerError,
			"Msg":    "新增失败",
			"Result": false,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code":   http.StatusOK,
			"Msg":    nil,
			"Result": succeed,
		})
	}
}

func DelNode(ctx *gin.Context) {
	name := ctx.Query("name")
	item := dao.ListItem{Name: name}
	if err := item.Delete(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"Code": http.StatusOK,
			"Msg":  nil,
		})
	}
}
