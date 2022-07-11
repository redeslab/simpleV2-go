package simple

import (
	"encoding/json"
	"fmt"
	"github.com/redeslab/go-simple/account"
	"github.com/redeslab/go-simple/network"
	"github.com/redeslab/go-simple/node"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	TXCheckTime   = time.Second * 2
	AndroidVerUrl = "https://redeslab.github.io/version.js"
	RuleDataUrl   = "https://redeslab.github.io/rule.txt"
	RuleVerUrl    = "https://redeslab.github.io/ruleVer.js"
)

type UICallBack interface {
	Log(str string)
	Notify(note string, data string)
	SysExit(err error)
}

type AppHelper struct {
	callback UICallBack
	timer    *time.Timer
	Wallet   account.Wallet
	ipCache  map[string]string
}

var appCaller = &AppHelper{
	timer:   time.NewTimer(TXCheckTime),
	ipCache: make(map[string]string),
}

func InitSystem(cb UICallBack) {
	appCaller.callback = cb
}

type PingResult struct {
	IP   string
	Ping float32
}

func minerIP(mid string) string {
	minerIP := appCaller.ipCache[mid]
	if len(minerIP) == 0 {
		minerIP = RefreshHostByAddr(mid)
		fmt.Println("======>>>find miner ip:=>", mid, minerIP)
	}

	return minerIP
}

func TestPing(mid string) []byte {

	minerIP := minerIP(mid)
	if minerIP == "" {
		return nil
	}

	mAddr := &net.UDPAddr{
		IP:   net.ParseIP(minerIP),
		Port: int(account.ID(mid).ToServerPort()),
	}
	timeOut := time.Second * 5

	fmt.Println("=====>start to ping:", minerIP)
	conn, err := net.DialTimeout("udp4", mAddr.String(), timeOut)
	if err != nil {
		fmt.Println("=====>dial miner err:", err)
		return nil
	}
	now := time.Now()
	defer conn.Close()
	testConn := network.JsonConn{Conn: conn}
	_ = testConn.SetDeadline(now.Add(timeOut))
	err = testConn.WriteJsonMsg(node.CtrlMsg{Typ: node.MsgPingTest, PT: &node.PingTest{
		PayLoad: mid,
	}})
	if err != nil {
		fmt.Println("=====>WriteJsonMsg err:", err)
		return nil
	}
	ctrlBuf := make([]byte, 1<<11)
	err = testConn.ReadJsonBuffer(ctrlBuf, &node.MsgAck{})
	if err != nil {
		fmt.Println("=====>ReadJsonMsg err:", err)
		return nil
	}

	result := PingResult{
		IP:   minerIP,
		Ping: float32(time.Now().Sub(now)) / float32(time.Millisecond),
	}
	fmt.Println("=====>finish to ping:", minerIP)
	data, _ := json.Marshal(result)
	return data
}

func MinerPort(addr string) int32 {
	mid := account.ID(addr)
	return int32(mid.ToServerPort())
}

func AndroidApkVersion() (ver string, err error) {
	data, err := getHttpJsonData(AndroidVerUrl)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type RuleVer struct {
	Ver int
}

func RuleVerInt() (int, error) {
	data, err := getHttpJsonData(RuleVerUrl)
	if err != nil {
		return -1, err
	}
	ver := &RuleVer{}
	if err := json.Unmarshal(data, ver); err != nil {
		return -1, err
	}

	return ver.Ver, nil
}

func RuleDataLoad() (string, error) {
	data, err := getHttpJsonData(RuleDataUrl)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getHttpJsonData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code is[%d]", resp.StatusCode)
		return nil, err
	}
	return body, nil
}
