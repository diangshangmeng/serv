package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func SendCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}

	err := service.SendCode(req.Phone)
	if err != nil {
		util.ResponseError(c, http.StatusInternalServerError, "发送失败")
		return
	}

	util.ResponseMessage(c, "发送成功")
}

func Register(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required,phone"`
		Password string `json:"password" binding:"required,min=6,max=32"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "参数错误：密码需要6-32位，包含大小写字母和数字")
		return
	}

	user, err := service.Register(req.Phone, req.Password)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"userInfo": gin.H{
			"id":         user.ID,
			"phone":      user.Phone,
			"status":     user.Status,
			"city_id":    user.CityID,
			"authStatus": user.AuthStatus,
		},
	}

	util.ResponseSuccess(c, result)
}

func Login(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required,phone"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}

	token, user, err := service.Login(req.Phone, req.Password)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"token": token,
		"userInfo": gin.H{
			"id":         user.ID,
			"phone":      user.Phone,
			"status":     user.Status,
			"city_id":    user.CityID,
			"authStatus": user.AuthStatus,
		},
	}

	util.ResponseSuccess(c, result)
}
