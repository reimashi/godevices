package serial

func (this *Config) GetPortAddress() string {
	return "/dev/tty" + this.port
}