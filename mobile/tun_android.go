//go:build android
// +build android

package simple

import (
	"errors"
	"fmt"
	"github.com/lightStarShip/go-tun2simple/stack"
	"github.com/lightStarShip/go-tun2simple/utils"
)

type ExtensionI interface {
	stack.TunDev
	stack.Wallet
	Log(s string)
	LoadRule() string
}

func InitEx(exi ExtensionI, logLevel int8) error {
	if exi == nil {
		return errors.New("invalid tun device")
	}
	utils.LogInst().InitParam(utils.LogLevel(logLevel), func(msg string, args ...any) {
		log := fmt.Sprintf(msg, args...)
		exi.Log(log)
	})
	rules := exi.LoadRule()
	return stack.Inst().SetupStack(exi, exi, rules)
}

func WritePackets(data []byte) (int, error) {
	return stack.Inst().WriteToStack(data)
}
