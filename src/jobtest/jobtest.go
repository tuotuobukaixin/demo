package main

import (
	"io"
	"net/http"
	"os"
	"crypto/sha256"
	"fmt"
)

var (
	url = "https://sm.myapp.com/original/game/TXsyzs-1.0.4204.123.exe"
)

func main() {

	res, err := http.Get(url)

	if err != nil {

		panic(err)

	}

	f, err := os.Create("TXsyzs-1.0.4204.123.exe")

	if err != nil {

		panic(err)

	}

	io.Copy(f, res.Body)
	file,err:=os.Open("TXsyzs-1.0.4204.123.exe")
	defer file.Close()
	hash:=sha256.New()
	io.Copy(hash,f)
	md:=hash.Sum(nil)
	sha256sum := fmt.Sprintf("%x", md)
	fmt.Println(sha256sum)
}
