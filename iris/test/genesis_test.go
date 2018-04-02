package test
//
//import (
//	"testing"
//	"io/ioutil"
//	"encoding/json"
//	"strings"
//	"os/exec"
//	"bytes"
//	"fmt"
//)
//
//type Coin struct {
//	Denom  string `json:"denom"`
//	Amount int64  `json:"amount"`
//}
//
//type Account struct {
//	Address string `json:"address"`
//	Coins   []Coin `json:"coins"`
//}
//
//type Options struct {
//	Accounts       []Account   `json:"accounts"`
//	Plugin_options interface{} `json:"plugin_options"`
//}
//
//type Genesis struct {
//	App_hash     string      `json:"app_hash"`
//	Chain_id     string      `json:"chain_id"`
//	Genesis_time string      `json:"genesis_time"`
//	Validators   interface{} `json:"validators"`
//	App_options  Options     `json:"app_options"`
//}
//
//func Test_modify(t *testing.T) {
//	path := "C:/Users/vincent/Desktop/"
//	bytes, err := ioutil.ReadFile(path + "genesis.json")
//	genesis := Genesis{}
//	if err == nil {
//		json.Unmarshal(bytes, &genesis)
//	}
//	account := genesis.App_options.Accounts
//	account1 := ""
//	if len(account) > 0 {
//		account1 = account[0].Address
//	}
//	accounts := AccountMap()
//	for _, v := range accounts {
//		if v != account1 {
//			account = append(account, newAccount(v))
//		}
//	}
//	genesis.App_options.Accounts = account
//	bytes, err = json.Marshal(genesis)
//	if err == nil {
//		ioutil.WriteFile(path+"genesisnew.json", bytes, 0644)
//	}
//}
//
//func newAccount(account string) Account {
//	coins := []Coin{{"iris", int64(100000)}}
//	return Account{account, coins}
//}
//
//func AccountMap() map[string]string {
//	c := exec.Command("iris", "client", "keys", "list")
//	var out bytes.Buffer
//	c.Stdout = &out
//	if err := c.Run(); err != nil {
//		fmt.Println("Error: ", err)
//	}
//	result := out.String()
//	s := strings.Split(result, "\n")
//	kv := make(map[string]string)
//	for i, v := range s {
//		if i == 0 {
//			continue
//		} else {
//			l := strings.Split(v, "\t\t")
//			if len(l) < 2 {
//				l = strings.Split(v, "\t")
//			}
//			if len(l) > 1 {
//				kv[l[0]] = l[1]
//			}
//		}
//	}
//	return kv
//}
