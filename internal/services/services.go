package services

import "pvz-service/config"

type Services struct {
	UserService      *UserService
	ProductService   *ProductService
	PvzService       *PVZService
	ReceptionService *ReceptionService
	Cfg              *config.Config
}
