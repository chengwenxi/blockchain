package test

import (
	"testing"
	"encoding/json"
	"container/list"
	"sync"
	"time"
	"math/rand"
	"github.com/chengwenxi/blockchain/iris/test/types"
	"github.com/chengwenxi/blockchain/iris/test/common"
)

var TO = "CAF62CF4258BB500D91C775106AD6419986B2A94"
var PASSWORD = "1234567890"
var SERVER = "http://116.62.62.39:1198"
var SERVERPOST = []string{"http://116.62.62.39:1198", "http://116.62.62.39:1298", "http://116.62.62.39:1398", "http://116.62.62.39:1498"}
var ch = make(chan int)
var signCh = make(chan int)
var l = list.New()
var lock sync.Mutex
var goNum = 200

func Test_PostTx(t *testing.T) {
	keys := getKeys()
	for i := 0; i < goNum; i++ {
		j := len(keys) / goNum
		go buildAndSignTxAll(keys, j*i, j)
	}
	for i := 0; i < goNum; i++ {
		<-signCh
	}
	println("sign end, number = ", l.Len())
	startTime := time.Now().Unix()
	for i := 0; i < goNum; i++ {
		go postTx()
	}
	for i := 0; i < goNum; i++ {
		<-ch
	}
	endTime := time.Now().Unix()
	println(endTime - startTime)
}

func buildAndSignTxAll(keys []types.Key, start int, number int) {
	for i := start; i < start+number; i++ {
		buildAndSignTx(keys[i].Name, keys[i].Address)
	}
	signCh <- 0
}

func postTx() {
	data := getPostTx()
	if data == nil {
		ch <- 0
		return
	}
	_ = common.DoPost(SERVERPOST[rand.Intn(3)]+"/tx", data)
	//println(string(result))
	postTx()
}

func getPostTx() []byte {
	defer lock.Unlock()
	lock.Lock()
	i1 := l.Back()
	var data []byte
	if i1 != nil {
		data = i1.Value.([]byte)
		l.Remove(i1)
	}
	return data
}

func getKeys() []types.Key {
	result := common.DoGet(SERVER + "/keys")
	var keys []types.Key
	if result != nil {
		json.Unmarshal(result, &keys)
	}
	return keys
}

func buildAndSignTx(name string, addr string) {

	amount := 1
	coin := "iris"
	result := common.DoGet(SERVER + "/query/nonce/" + addr)
	var nonce types.Nonce
	if result != nil {
		json.Unmarshal(result, &nonce)
	}
	nonce.Data += 1

	//build send
	si := new(types.SendInput)
	si.Amount = types.Coins{types.Coin{Denom: coin, Amount: int64(amount)}}
	si.From = &types.Actor{ChainID: "", App: "sigs", Address: addr}
	si.To = &types.Actor{ChainID: "", App: "sigs", Address: TO}
	//si.Fees = &Coin{Denom: feeCoin, Amount: 1}
	si.Sequence = nonce.Data
	siStr, _ := json.Marshal(si)
	result = common.DoPost(SERVER+"/build/send", siStr)
	if result == nil {
		return
	}

	//sign tx
	requestSign := new(types.RequestSign)
	requestSign.Name = name
	requestSign.Password = PASSWORD
	json.Unmarshal(result, &requestSign.Tx)
	rsStr, _ := json.Marshal(requestSign)
	result = common.DoPost(SERVER+"/sign", rsStr)
	if result == nil {
		return
	}
	//println(string(result))
	lock.Lock()
	l.PushBack(result)
	lock.Unlock()
	//send tx
}
