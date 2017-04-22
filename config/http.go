package config

const (
	//HostName - A default server only if os.Getenv("HOST") is not set
	HostName = "0.0.0.0"
	//HTTPPort - A default port only if os.Getenv("PORT") is not set
	HTTPPort = "443"
)
