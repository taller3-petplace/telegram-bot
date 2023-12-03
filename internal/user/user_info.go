package user

type UserInfo struct {
	id   int64
	name string
}

func NewUserInfo(id int64, name string) UserInfo {
	if name == "" {
		name = "Unknown"
	}

	return UserInfo{
		id:   id,
		name: name,
	}
}

func (ui UserInfo) GetName() string {
	return ui.name
}

func (ui UserInfo) GetUserID() int64 {
	return ui.id
}
