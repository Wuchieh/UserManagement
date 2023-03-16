package server

import (
	"github.com/Wuchieh/UserManagement/database"
	"log"
	"time"
)

var (
	logoutSession = make(chan string)
)

func cornTab() {
	runTime := 0
	for range time.Tick(10 * time.Second) {
		runTime++
		if runTime >= 60 { // 每10分鐘清除一次過期的Session
			runTime = 0
			go database.ReGetAccountCountDocuments()
			go database.ClearExpiredSession()
		}
		go func() {
			logoutSessionLise := func() (logoutSessions []string) {
				for {
					select {
					case s := <-logoutSession:
						logoutSessions = append(logoutSessions, s)
					default:
						return logoutSessions
					}
				}
			}()

			if len(logoutSessionLise) < 1 {
				return
			} else {
				err := database.ClearSessions(logoutSessionLise...)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}()

	}
}
