package simple

import (
	"flag"
	"fmt"
	"testing"
)

var (
	userAddr = ""
)

func init() {
	flag.StringVar(&userAddr, "uid", "", "--uid")
}

//go test -run  TestServerList
func TestServerList(t *testing.T) {
	ret := SyncServerList()
	fmt.Println("list:", string(ret))
}

//go test -run  TestQueryByAddr --uid
func TestQueryByAddr(t *testing.T) {
	ret := RefreshHostByAddr(userAddr)
	fmt.Println("ip:", string(ret))
}

func TestRuleVerInt(t *testing.T) {
	ret, err := RuleVerInt()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(ret))
}

func TestDnsRule(t *testing.T) {
	ret, err := RuleDataLoad()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ret)
}

func TestIpMust(t *testing.T) {
	ret, err := MustHitData()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ret)
}

func TestByPassDataLoad(t *testing.T) {
	ret, err := ByPassDataLoad()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ret)
}
