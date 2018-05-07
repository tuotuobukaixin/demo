package common

import (
	"testing"
)

func Test_ut_is_valid_ip(t *testing.T) {
	valid := false
	
	ip1 := "127.0.0.1"
	valid = Is_valid_ip(ip1)
	if valid != true {
		t.Error("fail here.")
		t.FailNow()
	}

	ip2 := "327.0.0.1"
	valid = Is_valid_ip(ip2)
	if valid != false {
		t.Error("fail here.")
		t.FailNow()
	}
	
}

func Test_ut_is_valid_ep(t *testing.T) {
	valid := false
	
	ep1 := "http//127.0.0.1"
	valid = Is_valid_endpoint(ep1)
	if valid != false {
		t.Error("fail here.")
		t.FailNow()
	}

	ep2 := "httpa://127.0.0.1"
	valid = Is_valid_endpoint(ep2)
	if valid != false {
		t.Error("fail here.")
		t.FailNow()
	}

	ep3 := "http://127.0.0.1:67000"
	valid = Is_valid_endpoint(ep3)
	if valid != false {
		t.Error("fail here.")
		t.FailNow()
	}
	
	ep10 := "http://127.0.0.1:5999"
	valid = Is_valid_endpoint(ep10)
	if valid != true {
		t.Error("fail here.")
		t.FailNow()
	}

	ep11 := "http://127.0.0.1"
	valid = Is_valid_endpoint(ep11)
	if valid != true {
		t.Error("fail here.")
		t.FailNow()
	}
}

