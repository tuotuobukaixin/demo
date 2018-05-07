package gcrypto

type Config struct {
	Path     string
	SyncPort string
}

var CryptoConf *Config

func init() {
	CryptoConf = &Config{
		SyncPort: "10143",
	}
}

func SetConfig(conf Config) {
	CryptoConf.SyncPort = conf.SyncPort
}
