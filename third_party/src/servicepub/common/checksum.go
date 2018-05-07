package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/astaxie/beego"
)

//每8KB取一次
const filechunk = 8192

func Md5sumFile(fullFileName string) string {
	file, err := os.Open(fullFileName)

	if err != nil {
		beego.Error("md5sum File failed, fullFileName:", fullFileName)
		return ""
	}

	defer file.Close()

	//获得文件大小
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(float64(filechunk), float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		//写入到hash中
		io.WriteString(hash, string(buf))
	}
	hashSum := fmt.Sprintf("%x", hash.Sum(nil))

	beego.Debug(fullFileName, "checksum is", hashSum)

	return hashSum
}
