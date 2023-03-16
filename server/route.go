package server

import (
	"github.com/gin-gonic/gin"
)

func routeInit(r *gin.Engine) {

	api := r.Group("/api")
	// User
	api.GET("/users/:page", getUsers) // 取得用戶列表
	//api.GET("/user/:userID", getUser)           // 取得用戶資訊
	api.GET("/user/userInfo", getUserInfo)      // 取得當前用戶資訊
	api.GET("/user/getUserCount", getUserCount) // 取得並刷新新用戶總量
	api.POST("/user", createUser)               // 創建
	api.PUT("/user", updateUser)                // 更新
	api.POST("/user/login", userLogin)          // 登入
	api.POST("/user/logout", userLogout)        // 登出
	api.POST("/user/register", userRegister)    // 註冊
	api.DELETE("/user", deleteUser)             // 刪除用戶(邏輯刪除, 若要使用物理刪除請使用其他工具)

}
