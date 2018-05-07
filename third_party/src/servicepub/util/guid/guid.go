package guid

import (
    "fmt"
    "math/rand"
    "os"
    "time"
guid "code.google.com/p/go-uuid/uuid"   
)

func Generate() string {
    return guid.New()
}

func Generate_old() string {
    b := createHex()
    return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func fileExist(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}

func createHex() []byte {
    result := make([]byte, 16)
    if fileExist("/dev/urandom") {
        f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
        f.Read(result)
        f.Close()
    } else {
        rand.Seed(time.Now().UTC().UnixNano())
        tmp := rand.Int63()
        rand.Seed(tmp)
        for i := 0; i < 16; i++ {
            result[i] = byte(rand.Intn(16))
        }
        result[6] = (result[6] & 0xF) | (4 << 4)
        result[8] = (result[8] | 0x40) & 0x7F

        //result[6] = (result[6] | 0x40) & 0x4F
        //result[8] = (result[8] | 0x80) & 0xBF
    }
    return result
}
