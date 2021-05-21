package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/CaptainLib"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getCaptainClient() *CaptainLib.CaptainClient {
	return CaptainLib.NewCaptainClient("http://192.168.1.224:5000/")
}

func getUrlIDParameter(name string, c *gin.Context) (int, error) {
	uriName := fmt.Sprintf("%sid", name)
	idVal, err := strconv.ParseInt(c.Param(uriName), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(idVal), nil
}

func forceIntRead(input string) int {
	idVal, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0
	}
	return int(idVal)
}