package script_service

import (
	"fmt"
	"github.com/MisakiOfScut/go-dage"
	"github.com/gin-gonic/gin"
	"go-dage-web/app/dao"
	"go.uber.org/zap"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
)

func GetAll(ctx *gin.Context) {
	name := ctx.Query("name")
	scripts := dao.Scripts{}
	if err := scripts.GetAll(name); err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "failed to get all versions of script " + name,
		})
	} else {
		ctx.JSON(200, gin.H{
			"Code":    http.StatusOK,
			"Msg":     nil,
			"Scripts": scripts,
		})
	}
}

func AddVersion(ctx *gin.Context) {
	script := dao.Script{}
	if err := ctx.BindJSON(&script); err != nil {
		ctx.JSON(500, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "",
		})
		return
	}

	if err := script.Add(); err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"Code": http.StatusOK,
			"Msg":  fmt.Sprintf("保存配置%s成功,版本号%d", script.Name, script.Version),
		})
	}
}

type mockGraphManager struct {
}

func (p *mockGraphManager) IsOprExisted(string2 string) bool {
	return true
}
func (p *mockGraphManager) GetOperatorInputs(oprName string) []string {
	return nil
}
func (p *mockGraphManager) GetOperatorOutputs(oprName string) []string {
	return nil
}
func (p *mockGraphManager) IsProduction() bool {
	return false
}

func Check(ctx *gin.Context) {
	script := dao.Script{}
	if err := ctx.BindJSON(&script); err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "发生内部错误，请联系管理员，错误信息：" + err.Error(),
		})
		return
	}

	dot, err := dage.TestBuildDAG(&script.Content)
	if err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusOK,
			"Msg":  err.Error(),
		})
		return
	}

	pngName := fmt.Sprintf("%s_%d",script.Name, rand.Int31())
	dotFile := fmt.Sprintf("web/static/temp/%s.dot", pngName)
	pngFile := fmt.Sprintf("web/static/temp/%s.png", pngName)
	webPath := fmt.Sprintf("temp/%s.png", pngName)

	err = ioutil.WriteFile(dotFile, []byte(dot), 0755)
	if nil != err {
		zap.S().Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "发生内部错误，请联系管理员，错误信息：" + err.Error(),
		})
		return
	}
	_, err = exec.Command("dot", "-Tpng", dotFile, "-o", pngFile).Output()
	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "发生内部错误，请联系管理员，错误信息：" + err.Error(),
		})
		return
	}

	zap.S().Infof("Write png into %s\n", pngFile)
	ctx.JSON(http.StatusOK, gin.H{
		"Code":    http.StatusOK,
		"Msg":     "",
		"PNGPath": webPath,
	})
}

func Publish(ctx *gin.Context) {
	script := dao.Script{}
	if err := ctx.BindJSON(&script); err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "发生内部错误，请联系管理员，错误信息：" + err.Error(),
		})
		return
	}

	if err := script.Get(); err != nil {
		ctx.JSON(200, gin.H{
			"Code": http.StatusInternalServerError,
			"Msg":  "发生内部错误，请联系管理员，错误信息：" + err.Error(),
		})
		return
	}

	if _, err := dage.TestBuildDAG(&script.Content); err != nil {
		ctx.JSON(200, gin.H{
			"Code": 201,
			"Msg":  err.Error(),
		})
		return
	}

	if err := script.Publish(); err != nil {
		ctx.JSON(200, gin.H{
			"Code": 201,
			"Msg":  err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"Code": 200,
		"Msg":  "",
	})
}