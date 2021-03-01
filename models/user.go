package models

type (
	User struct {
		Id       int    `json:"id"`
		User     string `json:"user"`
		Password string `json:"password"`
		Acc_type int    `json:"acc_type"`
	}
)
