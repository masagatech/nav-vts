package models

type VehicleModel struct {
	Id    int    `json:"id"`
	Imei  string `json:"imei"`
	Model string `json:"model"`
	Name  string `json:"name"`
}
