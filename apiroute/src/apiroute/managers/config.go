package managers

type Config struct {
	serviceDiscoveryClientAddr string
	listenPort                 string
}

func DefaultConfig() *Config {
	return &Config{
		serviceDiscoveryClientAddr: "http://127.0.0.1:10081",
		listenPort:                 "10080",
	}
}

func (c *Config) ServiceDiscoveryClientAddr() string {
	return c.serviceDiscoveryClientAddr
}

func (c *Config) SetServiceDiscoveryClientAddr(addr string) {
	c.serviceDiscoveryClientAddr = addr
}

func (c *Config) ListenPort() string {
	return c.listenPort
}

func (c *Config) SetListenPort(port string) {
	c.listenPort = port
}
