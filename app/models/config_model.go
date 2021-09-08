package models

type Config struct {
	Servers struct {
		Rest_server struct {
			Host string
			Port int
		}
		Tcp_server struct {
			Host string
			Port int
		}
	}
	Database struct {
		Database string
		Host     string
		Port     int
		User     string
		Pwd      string
	}
	Redis struct {
		Host string
		Port int
	}
	Tile38 struct {
		Host string
		Port int
	}
	Rabbitmq struct {
		Host string
		Port int
		User string
		Pwd  string
	}
}
