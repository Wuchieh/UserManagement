package server

import (
	"errors"
	"fmt"
	"github.com/Wuchieh/UserManagement/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// userRequestDate 接收用戶端發送過來的request body用
type userRequestDate struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Checked  bool   `json:"Checked"`
}

func (d userRequestDate) GetPasswordHash() string {
	var du database.Account
	du.Username = d.Username
	du.Password = d.Password
	return du.GetPasswordHash()
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func userRegisterLogic(c *gin.Context) (bool, error) { // 回傳是否為用戶輸入資料有誤 和錯誤理由
	var ra userRequestDate // register account data
	if err := c.Bind(&ra); err != nil {
		log.Println(err)
		return true, errors.New("資料解析錯誤")
	}

	if len(ra.Username) < 8 || len(ra.Password) < 8 {
		return true, errors.New("帳號或密碼長度不足 8")
	}

	u := database.NewAccount(ra.Username, ra.Password)
	if err := u.CreateAccount(); err != nil {
		if IsDup(err) {
			return true, errors.New("帳號已存在")
		} else {
			return false, errors.New("發生內部錯誤")
		}
	}
	return true, nil
}

func userLoginLogic(c *gin.Context) (bool, error) {
	var ra userRequestDate // register account data
	var uid string
	if err := c.Bind(&ra); err != nil {
		log.Println(err)
		return true, errors.New("資料解析錯誤")
	}

	if account, err := database.SearchAccount(ra.Username); err != nil {
		return false, errors.New("帳號不存在")
	} else {
		if ra.GetPasswordHash() != account.Password {
			return false, errors.New("密碼錯誤")
		} else {
			uid = account.GetID()
		}
	}

	session := sessions.Default(c)
	if ra.Checked {
		session.Options(
			sessions.Options{
				MaxAge: 604800, // 7天 單位 秒
				Path:   "/",
			})
	}

	t := 0
	u := uuid.New()
	var s *database.Session
	for {
		s = database.NewSession(u.String(), uid)
		err := s.CreateSession(
			func() int {
				if ra.Checked {
					return 604800
				}
				return 0
			}(),
		)
		if err == nil {
			break
		}
		if t > 5 {
			return false, errors.New("登入失敗 請新嘗試")
		}
		u = uuid.New()
		t++
	}
	session.Set("session", u.String())
	if err := session.Save(); err != nil {
		log.Println(err)
		if se, err := database.SearchSession(u.String()); err == nil {
			se.Drop()
		}
		return true, errors.New("發生內部錯誤")
	}
	return false, nil
}

func getUserList(limit, page int64) ([]*database.Account, error) {
	return database.GetAccountList(limit, page)
}

func userAuthentication(c *gin.Context) bool {
	session := sessions.Default(c)
	if a, err := database.GetAccountFromSession(fmt.Sprint(session.Get("session"))); err != nil {
		log.Println(err)
		return false
	} else {
		return a.Admin
	}
}

// 管理員修改用戶資訊
func adminUpdateUser(c *gin.Context) {
	ra := database.Account{}
	if err := c.Bind(&ra); err != nil {
		c.JSON(500, gin.H{"status": false, "msg": err.Error()})
		log.Println(err)
		return
	}
	ac, err := database.GetAccountFromID(ra.ID)
	if err != nil {
		c.JSON(500, gin.H{"status": false, "msg": err.Error()})
		log.Println(err)
		return
	}
	ac.Admin = ra.Admin
	if len(ra.Password) >= 8 {
		ac.Password = ra.Password
		ac.Password = ac.GetPasswordHash()
	}
	if err := ac.Update(); err != nil {
		c.JSON(500, gin.H{"status": false, "msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "更新成功"})
	}
}

func deleteUserLogic(c *gin.Context) {
	ra := database.Account{}
	if err := c.Bind(&ra); err != nil {
		c.JSON(500, gin.H{"status": false, "msg": err.Error()})
		log.Println(err)
		return
	}
	session := sessions.Default(c)
	if account, err := database.GetAccountFromSession(fmt.Sprint(session.Get("session"))); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	} else {
		if account.Username == ra.Username {
			c.JSON(400, gin.H{"status": false, "msg": "不可以刪除自己"})
			return
		} else if len(ra.Username) < 8 && len(ra.Password) < 8 {
			c.JSON(400, gin.H{"status": false, "msg": "資料錯誤"})
			return
		}
	}
	if err := ra.Delete(); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "成功刪除"})
	}
}

func createUserLogic(c *gin.Context) {
	var ra *database.Account
	if err := c.Bind(&ra); err != nil {
		log.Println(err)
		c.JSON(401, gin.H{"status": false, "msg": err.Error()})
		return
	}
	if ra.Admin {
		ra = database.NewAccount(ra.Username, ra.Password)
		ra.Admin = true
	} else {
		ra = database.NewAccount(ra.Username, ra.Password)
	}
	if err := ra.CreateAccount(); err != nil {
		log.Println(err)
		c.JSON(401, gin.H{"status": false, "msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "已成功創建新用戶"})
	}
}
