package engine

import (
	"fmt"

	"github.com/cpendery/cq/cq/engine/common"
	"github.com/cpendery/cq/cq/engine/postgres"
)

type Engine struct {
	driver common.Driver
	cache  map[string]common.DriverType
}

func NewEngine() *Engine {
	return &Engine{
		cache: loadCache(),
	}
}

func loadCache() map[string]common.DriverType {
	//cachePath := path.Join(xdg.CacheHome, internal.AppName)
	return make(map[string]common.DriverType)
}

func (e *Engine) getCachedType(username, dbname, host, port, password string) common.DriverType {
	return e.cache[username]
}

func (e *Engine) Connect(username, dbname, host, port, password, driverType string) error {
	if driverType != "" {
		dType := common.ToDriverType(driverType)
		for _, driver := range AllDrivers() {
			if driver.Type() == dType {
				e.driver = driver
				return e.driver.Connect(username, dbname, host, port, password)
			}
		}
		return fmt.Errorf("unknown driver: %s. Supported drivers are: %+v", driverType, common.AllDriverTypes)
	}

	cachedType := e.getCachedType(username, dbname, host, port, password)
	if cachedType != common.UnknownDriver {
		return fmt.Errorf("unimplemented cached type")
	}

	return fmt.Errorf("unimplemented no driver")
}

func AllDrivers() []common.Driver {
	return []common.Driver{
		postgres.NewPostgresDriver(),
	}
}
