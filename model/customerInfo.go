package model

type Customer struct {
	Name     string      `json:"Name"`
	Email    string      `json:"Email"`
	Address  string      `json:"Address"`
	JoinDate string      `json:"JoinDate"`
	Roles    interface{} `json:"Roles"`
}
