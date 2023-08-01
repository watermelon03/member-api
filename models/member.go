package models

type AdminAccount struct {
	AdminID           int    `json:"adminid"`
	AdminName         string `json:"adminname"`
	Password          string `json:"password"`
	RoleName          string `json:"rolename"`
	AdminState        string `json:"adminstate"`
	AccountUpdateDate string `json:"account-updatedate"`
}

type AdminInfo struct {
	AdminID        int    `json:"adminid"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Sex            string `json:"sex"`
	Birthday       string `json:"birthday"`
	Telephone      string `json:"tel"`
	Email          string `json:"email"`
	RegisterDate   string `json:"regisdate"`
	InfoUpdateDate string `json:"info-updatedate"`
}

type UserAccount struct {
	UserID            int    `json:"userid"`
	UserName          string `json:"username"`
	Password          string `json:"password"`
	LevelName         string `json:"levelname"`
	LevelImage        string `json:"levelimage"`
	UserState         string `json:"userstate"`
	UserPoint         string `json:"userpoint"`
	UserQR            string `json:"userqr"`
	AccountUpdateDate string `json:"account-updatedate"`
}

type UserInfo struct {
	UserID         int    `json:"userid"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Sex            string `json:"sex"`
	Birthday       string `json:"birthday"`
	Telephone      string `json:"tel"`
	Email          string `json:"email"`
	RegisterDate   string `json:"regisdate"`
	InfoUpdateDate string `json:"info-updatedate"`
	ImageName      string `json:"imagename"`
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterBody struct {
	UserID    int    `json:"userid"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	RoleID    int    `json:"roleid"`
	State     string `json:"state"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Sex       string `json:"sex" binding:"required"`
	Birthday  string `json:"birthday" binding:"required"`
	Telephone string `json:"tel" binding:"required"`
	Email     string `json:"email" binding:"required"`
}
