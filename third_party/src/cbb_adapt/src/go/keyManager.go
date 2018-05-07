package main

import (
	"bufio"
	"bytes"
	_ "cbb_adapt/src/go/aes"
	"cbb_adapt/src/go/gcrypto"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var LOGFILENAME string = "key_manager.log"

func WriteLog(logName string, logMsg string) {
	now := time.Now()
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	timestamp := fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d [paas] ",
		year, mon, day, hour, min, sec)
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

func WriteKeyFile(keyName string, key []byte) {
	keyFile, err := os.OpenFile(keyName, os.O_RDWR|os.O_CREATE, 0600)
	defer keyFile.Close()
	if err == nil {
		_, err = keyFile.Write(key)
	}
	if err != nil {
		WriteLog(LOGFILENAME, "write key error")
	}
	return
}

func ReadKeyFile(keyName string) ([]byte, error) {
	keyFile, err := os.Open(keyName)
	defer keyFile.Close()
	if nil == err {
		buf := make([]byte, 16)
		_, err := keyFile.Read(buf)
		if err != nil {
			log := fmt.Sprintf("read %s error", keyName)
			WriteLog(LOGFILENAME, log)
			return nil, err
		}

		return buf, nil
	}
	return nil, err
}

func SendKey(ip string, keyId string, buf []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(5 * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Second*3)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	}
	client := &http.Client{Transport: tr}

	requestUrl := "https://" + ip + "/key?keyId=" + keyId

	postBytesReader := bytes.NewReader(buf)
	rsp, err := client.Post(requestUrl, "text/plain", postBytesReader)
	if err != nil || rsp.StatusCode != 200 {

		log := fmt.Sprintf("post key to %s error", ip)
		fmt.Println(log)
		WriteLog(LOGFILENAME, log)
		return
	}
	fmt.Printf("rsp code is %d\n", rsp.StatusCode)
	log := fmt.Sprintf("post key to %s success", ip)
	fmt.Println(log)
	WriteLog(LOGFILENAME, log)
}

func ReadKeyId(keyIdFile string) (string, error) {
	keyIndexFile, err := os.Open(keyIdFile)
	defer keyIndexFile.Close()
	if nil == err {
		buf := make([]byte, 10)
		_, err := keyIndexFile.Read(buf)
		if err != nil {
			WriteLog(LOGFILENAME, "read keyId.txt error")
			return "", err
		}

		keyId := strings.Replace(string(buf), "\x00", "", -1)
		keyId = strings.Replace(keyId, "\n", "", -1)
		return keyId, nil
	}
	return "", err
}

func ReadIpSendKey(ipListFile string, keyId string, buf []byte) {
	ipFile, err := os.Open(ipListFile)
	defer ipFile.Close()
	if nil == err {
		buff := bufio.NewReader(ipFile)
		for {
			ip_s, err := buff.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			ip := strings.Replace(ip_s, "\n", "", -1)
			//fmt.Println(ip)
			SendKey(ip, keyId, buf)
		}
	} else {
		WriteLog(LOGFILENAME, "open ip.txt error")
	}

}

func GenerateKey() []byte {
	var list = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	var buf = make([]byte, 16)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		buf[i] = list[r.Intn(16)]
	}
	return buf
}

func main() {
	keyId, err := ReadKeyId("keyId.txt")
	if err != nil {
		WriteLog(LOGFILENAME, "read key id error")
		os.Exit(1)
	}
	fmt.Printf("keyId = %s\n", keyId)
	WriteLog(LOGFILENAME, "send key start")
	fmt.Println("send key start")
	var buf []byte
	eng, err := gcrypto.New("aes")
	if err != nil {
		fmt.Println("new eng err is: ", err)
		return
	}
	arg_num := len(os.Args)
	if arg_num == 1 {
		buf = GenerateKey()
		WriteKeyFile("key_"+keyId+".txt", buf)

		err := eng.FileEncrypt(0, "key_"+keyId+".txt", "key_"+keyId+"_encrypt.txt")
		if err != nil {
			fmt.Println("file encrypt error")
		}
		err = os.Remove("key_" + keyId + ".txt")
		if err != nil {
			fmt.Println(err)
		}
	} else if arg_num == 2 && os.Args[1] == "send" {
		err := eng.FileDecrypt(0, "key_"+keyId+"_encrypt.txt", "key_"+keyId+".txt")
		if err != nil {
			fmt.Println("file decrypt error")
		}
		buf, err = ReadKeyFile("key_" + keyId + ".txt")
		if err != nil {
			WriteLog(LOGFILENAME, "read key error")
			err = os.Remove("key_" + keyId + ".txt")
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(1)
		}
		err = os.Remove("key_" + keyId + ".txt")
		if err != nil {
			fmt.Println(err)
		}
	}

	ReadIpSendKey("ip.txt", keyId, buf)
	WriteLog(LOGFILENAME, "send key end")
	fmt.Println("send key end")
}
