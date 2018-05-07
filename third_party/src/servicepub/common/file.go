package common

import (
    "github.com/astaxie/beego"
      "os"
      "net/http"
      "io/ioutil"
      "path/filepath"
      "errors"
      "strconv"
      "strings"
      "time"
      "net"
)

type File interface {   
    Open(filepath string) (err error)
    Writestring(content string) (err error)
    StableWritestring(content string) (err error)   
    Size() (size int, err error)    
    Close()
    Rename(newname string) (err error)
    
    //funcs below can be used seperatelly
    File_exist(filepath string) (exist bool, err error) 
    Read(filepath string) (content []byte, err error)
    Read_remote(url string) (content []byte, err error)
    Read_ext(url string)(content []byte, err error)  //support http uri and os file path
    Dir_create(dir_abs_path string) (err error) 
}

type Dir interface {
    Get_entries(path string) (entries []string, err error)
}

func New_File() File{
    f := new(file_impl)

    return f
}

func New_Dir() Dir {
    d := new(dir_impl)

    return d
}

type file_impl struct {
    path string
    file *os.File
}

/*
    open file for reading and writing/appending, create it if not exist
*/
func (this *file_impl)Open(filepath string) (err error) {
    beego.Debug("Opening file: ", filepath)
    
    var exist = true;
    if _, err = os.Stat(filepath); os.IsNotExist(err) {
      exist = false;
    }

    if exist == true {      
        this.file, err = os.OpenFile(filepath, os.O_APPEND, 0660) 
    }else {
        beego.Debug("File not existed, will create file: ", filepath)
        this.file, err = os.Create(filepath)  
    } 

    this.path = filepath
    beego.Debug("file opened: ", filepath)
    return
}

func (this *file_impl)Writestring(content string) (err error){  
	beego.Debug("Write file, file name:", this.file.Name(), "content:", content)

	count, err := this.file.WriteString(content)
	if err != nil {
		beego.Warn("Write file fail, file name:", this.file.Name(), "error:", err)
	} else {
		beego.Debug("Write file successfully, file name:", this.file.Name(), "content length:", count)
	}
	
	return err
}

func (this *file_impl)StableWritestring(content string) (err error){  
    beego.Debug("stable writing file: ", this.file, "; ", content)
    
    count, err := this.file.WriteString(content) 
    if err != nil {
        beego.Warn("writing file fail: ", this.file, "; error:", err)
    } else {
        beego.Debug("writing file: ", this.file, "; ", count)
    }
    
    err = this.file.Sync()  
    if err != nil {
        beego.Warn("flush file fail: ", this.file, "; error:", err)
    } else {
        beego.Debug("flush file: ", this.file)
    }
    return 
}

func (this *file_impl)Size() (size int, err error) {
    return
}

func (this *file_impl)Close() {
    this.file.Close()
}

func (this *file_impl)Rename(newpath string) (err error){
    err = os.Rename(this.path, newpath)
    if err != nil {
        beego.Warn("Rename file fail: ", this.path, "; newpath: ", newpath, " ,error:", err)
    } else {
        beego.Debug("Rename file succeed: ", this.path, "; newpath: ", newpath)
    }

    return
}

func (this *file_impl)File_exist(filepath string) (exist bool, err error) {
    exist = true;
    if _, err = os.Stat(filepath); os.IsNotExist(err) {
      exist = false;
    }

    return
}

func (this *file_impl)Read(filepath string) (content []byte, err error){  
    beego.Debug("reading file: %s", filepath)

    exist := true
    exist, err = this.File_exist(filepath) 
    if ! exist || err != nil {
        err = EFileNotExist
        return
    }
    
    file,err := os.Open(filepath)  
    if err != nil {
        beego.Warn("open file error: ", filepath, ", ", err)
        return
    }  
    
    defer file.Close()     
    content, err = ioutil.ReadAll(file)
    return 
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

var itimeout int = 30
var timeout time.Duration = time.Duration(time.Duration(itimeout) * time.Second)  //second 
func dialTimeout(network, addr string) (net.Conn, error) {
        dt, err := net.DialTimeout(network, addr, timeout)
        if err != nil {
            beego.Warn("DialTimeout failed. error:", err)
            return dt, err
        }
        
        return dt, err
}
func (this *file_impl)Read_remote(url string) (contents []byte, err error) {
    var res *http.Response 

    tr := &http.Transport{
        Dial: dialTimeout,
    }
    client := &http.Client{
        Transport: tr,
        Timeout: timeout,
    }
        
    var req *http.Request
    req, err = http.NewRequest("GET", url, nil) 
    if err != nil {
        beego.Warn("new request failed. url:", url, ", error:", err)
        return
    }

    beego.Debug("Download remote file: ", url)
    res, err = client.Do(req)
    if res != nil {
        defer res.Body.Close()
    }
    
    if (err != nil)  ||  (res != nil && res.StatusCode != 200){
        beego.Warn("Download file failed. url:", url, ". error:", err)
        if (res != nil) {
            beego.Warn("Download file failed. error:", err, " ,statuscode:", res.StatusCode)    
        }
        if err == nil {
            err = errors.New("download file fail. statuscode: " + strconv.Itoa(res.StatusCode))
        }
        return
    }

    beego.Debug("Download file:", url, "from url success.")

    contents, err = ioutil.ReadAll(res.Body)
    if err != nil {
        beego.Warn("Read remote file fail. file:", url, ", err:", err)
        return
    }
    return contents, nil
}

func (this *file_impl)Read_ext(url string)(content []byte, err error) {
    beego.Debug("Reading file. url:", url)
    
    if strings.HasPrefix(url, "http://"){
        content, err = this.Read_remote(url)
        if err != nil {
            beego.Warn("Read remote file failed. url:", url, " ,error:", err)
            return
        }
    } else {
        content, err = this.Read(url)
        if err != nil {
            beego.Warn("Read local file failed, url: ", url , " ,error:", err)
            return
        }
    }

    return
}

func (this *file_impl)Dir_create(dir_abs_path string) (err error) {
    beego.Debug("Creating dir: ", dir_abs_path)
    
    exist := false
    exist,_ = this.File_exist(dir_abs_path)
    if exist == true {
        beego.Debug("Directory already exist. dir: ", dir_abs_path)
        return
    }
   
    err = os.MkdirAll(dir_abs_path, os.ModePerm)
    os.Chmod(dir_abs_path, 0750)
    if err != nil {
        beego.Warn("create fail. dir: ", dir_abs_path, ", error: ",err)
        return
    }

   return
}

type dir_impl struct {
    
}

func (this *dir_impl)Get_entries(path string) (entries []string, err error) {
    err = filepath.Walk(path, func(fpath string, f os.FileInfo, err error) error {
        if f == nil {
        return err
        }

        if f.IsDir() {
            beego.Debug("jump sub dir: ", fpath, ", error: ",err)
        } else  {
            beego.Debug("file found: ", fpath)
            entries = append(entries, fpath)
        }       
        
        return nil
    })
    
    if err != nil {
            beego.Debug("Get_entries returned %v\n", err)
    }

    return
}