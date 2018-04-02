package genesis

import (
	"os/exec"
	"bytes"
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/chengwenxi/blockchain/iris/test/types"
	"encoding/json"
)

func Modify(path string, denom string, amount int64) {
	bytes, err := ioutil.ReadFile(path + "genesis.json")
	genesis := types.Genesis{}
	if err == nil {
		json.Unmarshal(bytes, &genesis)
	}
	account := genesis.App_options.Accounts
	account1 := ""
	if len(account) > 0 {
		account1 = account[0].Address
	}
	accounts := AccountMap()
	for _, v := range accounts {
		if v != account1 {
			account = append(account, newAccount(v, denom, amount))
		}
	}
	genesis.App_options.Accounts = account
	bytes, err = json.Marshal(genesis)
	if err == nil {
		ioutil.WriteFile(path+"genesis.json", bytes, 0644)
	}
}

func newAccount(addr string, denom string, amount int64) types.Account {
	coins := []types.Coin{{denom, amount}}
	return types.Account{addr, coins}
}

func AccountMap() map[string]string {
	c := exec.Command("iris", "client", "keys", "list")
	var out bytes.Buffer
	c.Stdout = &out
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	result := out.String()
	s := strings.Split(result, "\n")
	kv := make(map[string]string)
	for i, v := range s {
		if i == 0 {
			continue
		} else {
			l := strings.Split(v, "\t\t")
			if len(l) < 2 {
				l = strings.Split(v, "\t")
			}
			if len(l) > 1 {
				kv[l[0]] = l[1]
			}
		}
	}
	return kv
}
