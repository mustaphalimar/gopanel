package database

type Database interface {
	Close()
	Ping() error
	Health() error
}
