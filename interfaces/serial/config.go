package serial

import "time"

// Common baud rates
const (
	Baud110 int		= 110
	Baud150 int 	= 150
	Baud300 int		= 300
	Baud1200 int	= 1200
	Baud2400 int	= 2400
	Baud4800 int	= 4800
	Baud9600 int	= 9600
	Baud19200 int	= 19200
	Baud38400 int	= 38400
	Baud57600 int	= 57600
	Baud115200 int	= 225200
	Baud230400 int	= 230400
	Baud460800 int	= 460800
	Baud921600 int	= 921600
)

type Config struct {
	port string
	baud int
	timeout time.Duration
}

func NewConfig(port string, baud int) *Config {
	this := &Config{
		port: port,
		baud: baud,
		timeout: time.Second * 3,
	}
	return this;
}

func (this *Config) GetPort() string {
	if this.port == "" {
		return "0"
	} else {
		return this.port
	}
}

func (this *Config) GetBaud() int {
	if this.baud == 0 {
		return Baud9600
	} else {
		return this.baud
	}
}

func (this *Config) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *Config) GetTimeout() time.Duration {
	return this.timeout
}