package server

import (
	"fmt"
	"github.com/Wuchieh/UserManagement/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func deleteUser(c *gin.Context) {
	if userAuthentication(c) {
		deleteUserLogic(c)
	}
}

func updateUser(c *gin.Context) {
	if userAuthentication(c) {
		adminUpdateUser(c)
	} else {
		c.JSON(400, gin.H{"status": false, "msg": "權限不足"})
	}
}

// 需進行身分驗證
func createUser(c *gin.Context) {
	if userAuthentication(c) {
		createUserLogic(c)
	}
}

func getUserCount(c *gin.Context) {
	if userAuthentication(c) {
		c.String(200, strconv.Itoa(int(database.ReGetAccountCountDocuments())))
	} else {
		c.String(404, "error")
	}
}

func getUsers(c *gin.Context) {
	session := sessions.Default(c)

	limitAtoi, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limitAtoi = 10
	}
	limit := int64(limitAtoi)

	atoi, err := strconv.Atoi(c.Param("page"))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"status": false, "msg": "攔截異常輸入"})
		return
	}

	if a, err := database.GetAccountFromSession(fmt.Sprint(session.Get("session"))); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"status": false, "msg": "管理員資料取得失敗"})
	} else if !a.Admin {
		c.JSON(404, gin.H{"status": false, "msg": "管理員身分驗證失敗"})
	} else {
		if accounts, err := getUserList(limit, int64(atoi)); err != nil {
			c.JSON(500, gin.H{"status": false, "msg": "用戶資料取得失敗"})
		} else {
			c.JSON(200, gin.H{
				"status":         true,
				"userList":       accounts,
				"countDocuments": database.GetAccountCountDocuments(),
			})
		}
	}

}

//func getUser(c *gin.Context) {
//
//}

func getUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	if a, err := database.GetAccountFromSession(fmt.Sprint(session.Get("session"))); err != nil {
		log.Println(err)
		session.Options(sessions.Options{Path: "/", MaxAge: -1})
		if err := session.Save(); err != nil {
			log.Println(err)
		}
		c.JSON(404, gin.H{
			"status": false,
			"msg":    err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"status":   true,
			"Username": a.Username,
			"Admin":    a.Admin,
		})
	}
}

func userLogin(c *gin.Context) {
	if b, err := userLoginLogic(c); err != nil {
		if b {
			c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		} else {
			c.JSON(500, gin.H{"status": false, "msg": err.Error()})
		}
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "登入成功"})
	}
}

func userRegister(c *gin.Context) {
	if b, err := userRegisterLogic(c); err != nil {
		if b {
			c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		} else {
			c.JSON(500, gin.H{"status": false, "msg": err.Error()})
		}
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "註冊成功"})
	}
}

func userLogout(c *gin.Context) {
	// 直接將前端的session設置為過期
	session := sessions.Default(c)
	dbSession := session.Get("session")
	session.Options(sessions.Options{MaxAge: -1, Path: "/"})
	if err := session.Save(); err != nil {
		log.Println(err)
	}
	c.String(200, "已登出")

	// 後端處理部分 主要目的為清除db中的session
	if dbSession != nil {
		go func() {
			a := fmt.Sprint(dbSession)
			logoutSession <- a
		}()
	}
}
