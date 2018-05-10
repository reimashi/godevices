package serial

func (this *Config) GetPortAddress() string {
	return "COM" + this.port
}