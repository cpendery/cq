package engine

import (
	"fmt"
	"sync"

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
	return common.UnknownDriver
	//e.cache[username]
}

func (e *Engine) Connect(username, dbname, host, port, password, driverType string) error {
	//use the given type specified
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

	//try to use the cached connection for this connection
	cachedType := e.getCachedType(username, dbname, host, port, password)
	if cachedType != common.UnknownDriver {
		return fmt.Errorf("unimplemented cached type")
	}

	//try to connect using every driver and continue as soon as a connection is made with 1
	var (
		newDriver common.Driver = nil
		setLock   sync.RWMutex
		doneLock  sync.RWMutex
		doneOnce  sync.Once
		wg        sync.WaitGroup
		set       bool = false
	)
	drivers := AllDrivers()
	doneLock.Lock()
	wg.Add(len(drivers))
	for _, d := range drivers {
		driver := d
		go func() {
			defer wg.Done()
			err := driver.Connect(username, dbname, host, port, password)
			setLock.Lock()
			switch {
			case !set && err == nil:
				newDriver = driver
				set = true
				setLock.Unlock()
				doneOnce.Do(func() { doneLock.Unlock() })
			case err == nil:
				driver.Close()
				setLock.Unlock()
			default:
				setLock.Unlock()
			}
		}()
	}
	go func() {
		wg.Wait()
		doneOnce.Do(func() { doneLock.Unlock() })
	}()

	doneLock.Lock()
	e.driver = newDriver
	if newDriver == nil {
		return fmt.Errorf("unable to connect to a database of any type with the given parameters")
	}
	return nil
}

func AllDrivers() []common.Driver {
	return []common.Driver{
		postgres.NewPostgresDriver(),
	}
}
