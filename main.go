package main

import (
	"encoding/json"
	"github.com/Wuchieh/UserManagement/database"
	"github.com/Wuchieh/UserManagement/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type setting struct {
	SessionStore       string `json:"sessionStore"`
	Ip                 string `json:"ip"`
	Port               string `json:"port"`
	SessionExpiredTime int    `json:"sessionExpiredTime"`
	GinMode            string `json:"ginMode"`
}

var (
	s    setting
	file []byte
	err  error
)

func main() {
	sc := make(chan os.Signal, 1)

	if file, err = os.ReadFile("setting.json"); err != nil {
		log.Panic(err)
	}
	if err = json.Unmarshal(file, &s); err != nil {
		log.Panic(err)
	}
	go func() {
		ip := s.Ip + ":" + s.Port
		if err := server.Run(ip, s.SessionStore, s.SessionExpiredTime, s.GinMode); err != nil {
			log.Println(err)
		}
	}()

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer func() {
		database.Disconnect()
		if err = server.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
}
