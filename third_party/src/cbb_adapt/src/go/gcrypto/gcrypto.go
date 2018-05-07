package gcrypto

import (
	"fmt"
)

type Engine interface {
	Encrypt(domainId int, data string) (string, error)
	Decrypt(domainId int, encData string) (string, error)
	FileEncrypt(domainId int, srcPath string, dstPath string) error
	FileDecrypt(domainId int, srcPath string, dstPath string) error
}

type InitFunc func() (Engine, error)

var (
	engines map[string]InitFunc
)

func init() {
	engines = make(map[string]InitFunc)
}

func Register(algo string, initFunc InitFunc) error {
	if _, exists := engines[algo]; exists {
		return fmt.Errorf("engine already registered %s", algo)
	}

	engines[algo] = initFunc

	return nil
}

func New(algo string) (engine Engine, err error) {
	if initFunc, exists := engines[algo]; exists {
		return initFunc()
	}

	return nil, fmt.Errorf("no such algorithm %s", algo)
}
