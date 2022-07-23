package common

type Driver interface {
	Connect(username, dbname, host, port, password string) error
	Type() DriverType
}

type DriverType string

const (
	PostgresDriver = "postgres"
	UnknownDriver  = "unknown"
)

var AllDriverTypes = []DriverType{
	PostgresDriver,
}

func ToDriverType(driver string) DriverType {
	switch driver {
	case PostgresDriver:
		return PostgresDriver
	default:
		return UnknownDriver
	}
}
