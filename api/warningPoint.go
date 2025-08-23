package main

import (
	"net/http"
	_ "template-mobile-app-api/docs"

	"github.com/gin-gonic/gin"
)

func GetWarningPoints(c *gin.Context) {
	//ここで結合する
	warningPoints := worningIntersectionPoints

	c.JSON(http.StatusOK, warningPoints)
}
