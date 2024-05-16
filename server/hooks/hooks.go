package hooks

import "server/config"

func New() *Service {
	return &Service{
		Cfg: *config.Get(),
	}
}
