package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cpendery/cq/cq/engine/common"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	name common.DriverType
}

func NewPostgresDriver() common.Driver {
	return &PostgresDriver{
		name: "postgres",
	}
}

func (p PostgresDriver) Type() common.DriverType {
	return p.name
}

func getDataSourceName(username, dbname, host, port, password, sslmode string) string {
	var sb strings.Builder

	if username != "" {
		sb.WriteString(fmt.Sprintf("user=%s ", username))
	}
	if dbname != "" {
		sb.WriteString(fmt.Sprintf("dbname=%s ", dbname))
	}
	if host != "" {
		sb.WriteString(fmt.Sprintf("host=%s ", host))
	}
	if port != "" {
		sb.WriteString(fmt.Sprintf("port=%s ", port))
	}
	if password != "" {
		sb.WriteString(fmt.Sprintf("password=%s ", password))
	}
	switch sslmode {
	case "":
		sb.WriteString("sslmode=disable")
	default:
		sb.WriteString(fmt.Sprintf("sslmode=%s ", password))
	}

	return sb.String()
}

func (p *PostgresDriver) Connect(username, dbname, host, port, password string) error {
	dataSource := getDataSourceName(username, dbname, host, port, password, "")
	db, err := sql.Open(string(p.name), dataSource)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")
	return nil
}

func (p *PostgresDriver) Close() error {
	return nil
}
