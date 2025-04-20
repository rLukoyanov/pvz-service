package services

import (
	"pvz-service/config"
)

type Services struct {
	UserService      UserServiceInterface
	ProductService   ProductServiceInterface
	PvzService       PVZServiceInterface
	ReceptionService ReceptionServiceInterface
	Cfg              *config.Config
}
