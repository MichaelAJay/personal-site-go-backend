package config

type Config struct {
	Env   string `json:"env"`
	DbDsn string `json:"db_dsn"`
}
