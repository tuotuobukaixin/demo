package aes

/*
#cgo LDFLAGS: -L ../../../../../../build/crypto/output/ -laes_crypto -ldl
#include <stdlib.h>
#include "aes_crypto.h"
*/
import "C"

import (
	"cbb_adapt/src/go/gcrypto"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	interName                = "eth0"
	START_KEY_SERVER_SUCCESS = 0
	START_KEY_SERVER_FAILURE = 1
)

type Engine struct {
	name string
}

var registerFlag = 0
var LOGFILENAME string = "/var/paas/paas_crypto_log.log"

func init() {
	WriteLog(LOGFILENAME, "init aes start...")
	gcrypto.Register("aes", Init)
}

func WriteLog(logName string, logMsg string) {
	now := time.Now()
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	timestamp := fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d [paas][%d] ",
		year, mon, day, hour, min, sec, os.Getpid())
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0640)
	defer logFile.Close()
	if err == nil {
		_, err = logFile.WriteString(timestamp + logMsg + "\n")
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}

func ReadCryptoSwitch(cryptoFileName string) (string, error) {
	cryptoFile, err := os.Open(cryptoFileName)
	defer cryptoFile.Close()
	if nil == err {
		buf := make([]byte, 10)
		_, err := cryptoFile.Read(buf)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		cryptoFlag := strings.Replace(string(buf), "\x00", "", -1)
		cryptoFlag = strings.Replace(string(cryptoFlag), "\n", "", -1)
		return cryptoFlag, nil
	}
	return "", err
}

func handlerKey(rsp http.ResponseWriter, req *http.Request) {
	keyIdString := req.FormValue("keyId")
	keyId, err := strconv.Atoi(keyIdString)
	if err != nil {
		fmt.Println(err)
		WriteLog(LOGFILENAME, "Change keyId String to int error")
		return
	}
	if strings.ToUpper(req.Method) == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			WriteLog(LOGFILENAME, "read body error")
			return
		}
		key := string(body)
		log := fmt.Sprintf("try to register key, Id is %d", keyId)
		WriteLog(LOGFILENAME, log)

		//如果不是第一个分发的密钥，则将之前的密钥设置成不可加密，但是可以解密
		if keyId > 0 {
			err = SetKeyInvalid(1, keyId-1)
			if err != nil {
				WriteLog(LOGFILENAME, "set key invalid error")
				return
			}
		}
		err = RegisterKey(1, keyId, key, len(key))
		if err != nil {
			WriteLog(LOGFILENAME, "register Key error")
			rsp.WriteHeader(405)
			rsp.Write([]byte("register Key error"))
			return
		}
		WriteLog(LOGFILENAME, "register Key ok")
		return
	}
	log := fmt.Sprintf("request method is %s", req.Method)
	WriteLog(LOGFILENAME, log)
}

func KeyServer(ch chan int) {
	WriteLog(LOGFILENAME, "enter KeyServer()")
	http.HandleFunc("/key", handlerKey)
	ifi, err := net.InterfaceByName(interName)
	if err != nil {
		fmt.Println(err)
		WriteLog(LOGFILENAME, "get interface error")
		ch <- START_KEY_SERVER_FAILURE
		return
	}
	WriteLog(LOGFILENAME, "get interface successfully")
	addrs, err := ifi.Addrs()
	if err != nil {
		fmt.Println(err)
		WriteLog(LOGFILENAME, "get addrs error")
		ch <- START_KEY_SERVER_FAILURE
		return
	}
	WriteLog(LOGFILENAME, "get addrs successfully")
	ips := strings.Split(addrs[0].String(), "/")
	WriteLog(LOGFILENAME, "begin to start server")
	ch <- START_KEY_SERVER_SUCCESS

	log := fmt.Sprintf("start gcrypto server port %s", gcrypto.CryptoConf.SyncPort)
	WriteLog(LOGFILENAME, log)
	err = http.ListenAndServeTLS(ips[0]+":"+gcrypto.CryptoConf.SyncPort, "cert.pem", "key.pem", nil)
	if err != nil {
		WriteLog(LOGFILENAME, fmt.Sprintf("Failed to start key server, error occurs: %s.", err.Error()))
	}
}

func checkNeedSharedCrypto(cryptoFileName string) bool {
	_, err := os.Stat(cryptoFileName)
	if err == nil {
		cryptoFlag, err := ReadCryptoSwitch(cryptoFileName)
		if err == nil && cryptoFlag == "true" {
			return true
		}
	}
	return false
}

func Init() (gcrypto.Engine, error) {
	WriteLog(LOGFILENAME, "Enter aes.Init()")
	if ret := C.aesInit(); ret != C.RET_SUCCESS {
		return nil, fmt.Errorf("engine init failed, return %d", ret)
	}
	WriteLog(LOGFILENAME, "Engine init success")
	eng := &Engine{
		name: "aes",
	}

	//只起一次KeyServer
	if registerFlag == 0 && checkNeedSharedCrypto("crypto_switch.txt") {
		WriteLog(LOGFILENAME, "Try to start keyServer")
		registerFlag = 1
		eng, err := gcrypto.New("aes")
		if err != nil {
			WriteLog(LOGFILENAME, " new eng err is : "+err.Error())
			return eng, nil
		}
		_, err = os.Stat("key.pem")
		if err != nil {
			WriteLog(LOGFILENAME, "There is no private key file")
			_, err = os.Stat("key_encrypt.pem")
			if err == nil {
				err = eng.FileDecrypt(0, "key_encrypt.pem", "key.pem")
				if err != nil {
					WriteLog(LOGFILENAME, "decrypt private key file failed")
					return eng, nil
				}
				err := os.Chmod("key_encrypt.pem", 0600)
				if err != nil {
					WriteLog(LOGFILENAME, "chmod private key file failed")
				}
			} else {
				WriteLog(LOGFILENAME, "no encrypt key file error")
				return eng, nil
			}
		} else {
			WriteLog(LOGFILENAME, "There is a private key file")
		}
		ch := make(chan int)
		go KeyServer(ch)
		<-ch
		_, err = os.Stat("key_encrypt.pem")
		if err != nil {
			err = eng.FileEncrypt(0, "key.pem", "key_encrypt.pem")
			if err != nil {
				WriteLog(LOGFILENAME, "encrypt private key file failed")
				return eng, nil
			}
			WriteLog(LOGFILENAME, "encrypt private key file successfully")
			err := os.Chmod("key_encrypt.pem", 0600)
			if err != nil {
				WriteLog(LOGFILENAME, "chmod private key file failed")
			} else {
				WriteLog(LOGFILENAME, "chmod private key file successfully")
			}
		}
		time.Sleep(time.Duration(5000) * time.Millisecond)
		err = os.Remove("key.pem")
		if err != nil {
			WriteLog(LOGFILENAME, "failed to remove key file")
		} else {
			WriteLog(LOGFILENAME, "remove key file successfully")
		}
	} else {
		if registerFlag == 1 {
			WriteLog(LOGFILENAME, "It already has tried to start, and this time no need to do so.")
		} else {
			WriteLog(LOGFILENAME, "There is no need to share, so it no need to start.")
		}
	}
	WriteLog(LOGFILENAME, "Leave aes.Init()")
	return eng, nil
}

func RegisterKey(domainId int, keyId int, key string, keyLen int) error {
	if ret := C.registerWorkingKey(C.int(domainId), C.int(keyId), C.CString(key), C.int(keyLen)); ret != C.RET_SUCCESS {
		return fmt.Errorf("register key error")
	}
	return nil
}

func SetKeyInvalid(domainId int, keyId int) error {
	if ret := C.setKeyInvalid(C.int(domainId), C.int(keyId)); ret != C.RET_SUCCESS {
		return fmt.Errorf("set key invalid error")
	}
	return nil
}

func (Engine) Encrypt(domainId int, data string) (string, error) {
	var encdata *C.char
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	if encdata = C.aesEncrypt(C.int(domainId), cdata); encdata == nil {
		return "", fmt.Errorf("encrypt error")
	}
	defer C.free(unsafe.Pointer(encdata))
	return C.GoString(encdata), nil
}

func (Engine) Decrypt(domainId int, data string) (string, error) {
	var decdata *C.char
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	if decdata = C.aesDecrypt(C.int(domainId), cdata); decdata == nil {
		return "", fmt.Errorf("encrypt error")
	}
	defer C.free(unsafe.Pointer(decdata))

	return C.GoString(decdata), nil
}

func (Engine) FileEncrypt(domainId int, srcPath string, dstPath string) error {

	cSrcPath := C.CString(srcPath)
	cDstPath := C.CString(dstPath)

	defer C.free(unsafe.Pointer(cSrcPath))
	defer C.free(unsafe.Pointer(cDstPath))

	if ret := C.aesFileEncrypt(C.int(domainId), cSrcPath, cDstPath); ret != C.RET_SUCCESS {
		return fmt.Errorf("file encrypt error")
	}

	return nil
}

func (Engine) FileDecrypt(domainId int, srcPath string, dstPath string) error {

	cSrcPath := C.CString(srcPath)
	cDstPath := C.CString(dstPath)

	defer C.free(unsafe.Pointer(cSrcPath))
	defer C.free(unsafe.Pointer(cDstPath))

	if ret := C.aesFileDecrypt(C.int(domainId), cSrcPath, cDstPath); ret != C.RET_SUCCESS {
		return fmt.Errorf("file decrypt error")
	}

	return nil
}
