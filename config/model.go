package config

type Config struct {
	Elastic ElasticsearchConfig `mapstructure:"elasticsearch"`
}

type ElasticsearchConfig struct {
	Host                   string `mapstructure:"host"`
	Username               string `mapstructure:"username"`
	Password               string `mapstructure:"password"`
	SolarIndex             string `mapstructure:"solar_index"`
	DailyProductionIndex   string `mapstructure:"daily_production_index"`
	MonthlyProductionIndex string `mapstructure:"monthly_production_index"`
}
