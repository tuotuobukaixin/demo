package common

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

func ZipFile(fullFileName string, destPath string) (string, error) {
	dir, fileNameNoSuffix, fileSuffix := SplitFullFileName(fullFileName)

	//如果目标目录为空，则放在文件所在目录
	if destPath == "" {
		destPath = dir
	}

	var zipFileName string = destPath + "/" + fileNameNoSuffix + ".zip"

	//return "", os.ErrInvalid

	fw, err := os.Create(zipFileName)

	//如果出错，则删除zip文件
	defer func() {
		fw.Close()
		if err != nil {
			os.Remove(zipFileName)
		}
	}()

	if err != nil {
		beego.Error("Create zip file failed, error:", err)
		return "", err
	}

	w := zip.NewWriter(fw)
	defer w.Close()

	srcFile, err := os.Open(fullFileName)
	defer srcFile.Close()
	if err != nil {
		beego.Error("Open UDR file failed, error:", err)
		return "", err
	}

	header := &zip.FileHeader{Name: fileNameNoSuffix + fileSuffix}
	header.Method = zip.Deflate
	
	//由于SetModTime中会转化成UTC时间，导致文件修改时间显示不对，所以直接将当前时间的字面时间作为UTC时间
	header.SetModTime(TimeLiteralToUTC(time.Now()))

	destFile, err := w.CreateHeader(header)
	//defer destFile.
	if err != nil {
		beego.Error("Create file in zip failed, error:", err)
		return "", err
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		beego.Error("Write zip file failed, error:", err)
		return "", err
	}

	err = w.Flush()

	if err != nil {
		beego.Error("Flush zip file failed, error:", err)
		return "", err
	}

	beego.Debug("Zip file success, fullFileName:", fullFileName, "zipFileName:", zipFileName)

	return zipFileName, nil
}

//折分全路径文件名，分成路径，文件名称（无后缀），文件后缀
func SplitFullFileName(fullFileName string) (string, string, string) {
	dir, fileNameWithSuffix := path.Split(fullFileName)
	dir = strings.TrimRight(strings.TrimRight(dir, "/"), "\\")

	fileSuffix := path.Ext(fileNameWithSuffix)
	fileNameNoSuffix := strings.TrimSuffix(fileNameWithSuffix, fileSuffix)

	return dir, fileNameNoSuffix, fileSuffix
}

//将时间的字面值转换成UTC时间，即如果是东8区，则相当于8个小时
func TimeLiteralToUTC(t time.Time) time.Time {
	location, _ := time.LoadLocation("UTC")
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return time.Date(year, month, day, hour, minute, second, 0, location)
}
