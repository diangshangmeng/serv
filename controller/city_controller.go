package controller

import (
	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func GetCityList(c *gin.Context) {
	cities, err := service.GetCityList()
	if err != nil {
		util.ResponseError(c, 500, "获取城市列表失败")
		return
	}

	util.ResponseSuccess(c, cities)
}
