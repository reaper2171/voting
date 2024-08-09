package models

type RegReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutReq struct {
	Username string `json:"username"`
}
