package service

import (
	"com.fs/event-service/db"
)

// Init 初始化service
func Init() {
	db.Open()
}

// Destroy 销毁service
func Destroy() {
	db.Close()
}
