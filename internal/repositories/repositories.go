package repositories

import "pvz-service/config"

type Repos struct {
	AuthRepo      *UserRepository
	ProductRepo   *ProductRepository
	PvzRepo       *PVZRepository
	ReceptionRepo *ReceptionRepository
	Cfg           *config.Config
}
