package main

import (
	"net"
	"fmt"
	"os"
	"flag"
	"time"
	"strings"
	"strconv"
)

func dnstest(i int, hostname string, ipstring string)  {
	start := time.Now()
	index := 0
	for {
		index++
		ips := strings.Split(ipstring, ",")
		ns, err := net.LookupHost(hostname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Err: %s", err.Error())
			continue
		}
		flg := false
		for _,ip :=range ips{
			if  ns[0] == ip {
				flg = true
				break
			}
		}
		if  !flg {
			fmt.Fprintln(os.Stdout, "ip err:ip: "+ipstring+ "nsip: "+ ns[0])
		}
		now := time.Now()
		subM := now.Sub(start)

		if int(subM.Seconds()) >= 1 {
			start = time.Now()
			fmt.Println(strconv.Itoa(index)+" " + strconv.Itoa(int(subM.Seconds())))
			index = 0
		}
	}
	return
}

func main() {
	num := flag.Int("num", 10, "theard num")
	hostname := flag.String("hostname", "", "hostname")
	ip := flag.String("ip", "", "ip")
	flag.Parse()
	fmt.Println(*hostname + *ip)
	for a := 0; a < *num; a++ {

		go dnstest(a, *hostname, *ip)
	}
	time.Sleep(360000 * time.Hour)
}
