package db

import (
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/user"
)

type FakeDB struct {
	db map[int64]user.UserInfo
}

func NewFakeDB() *FakeDB {
	db := make(map[int64]user.UserInfo)
	return &FakeDB{
		db: db,
	}
}

func (fake *FakeDB) AddUser(telegramUserData tele.User) {
	userInfo := user.NewUserInfo(telegramUserData.ID, telegramUserData.FirstName)
	fake.db[telegramUserData.ID] = userInfo
}

func (fake *FakeDB) GetUser(userID int64) (user.UserInfo, bool) {
	userInfo, found := fake.db[userID]
	return userInfo, found
}
