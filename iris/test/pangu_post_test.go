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
	"strconv"
)

var TO = "CAF62CF4258BB500D91C775106AD6419986B2A94"
var PASSWORD = "1234567890"
var SERVER = "http://116.62.62.39:1198"
//var SERVERPOSTS = []string{"http://116.62.62.39:1198", "http://116.62.62.39:1298", "http://116.62.62.39:1398", "http://116.62.62.39:1498"}
var SERVERPOSTS = []string{"http://116.62.62.39:1198"}
var ch = make(chan int)
var signCh = make(chan int)
var l = list.New()
var lock sync.Mutex
var goNum = 1
var goGetNum = 10
var minPostTime = int64(10000000000)
var maxPostTime = int64(0)
var successNum = int64(0)
var totalTime = int64(0)

func Test_PostTx(t *testing.T) {
	keys := getKeys()
	keys = keys[0:10]
	for i := 0; i < goGetNum; i++ {
		j := len(keys) / goGetNum
		go buildAndSignTxAll(keys, j*i, j)
	}
	for i := 0; i < goGetNum; i++ {
		<-signCh
	}
	println("sign end, number = ", l.Len())
	for i := 0; i < goNum; i++ {
		go postTx()
	}
	for i := 0; i < goNum; i++ {
		<-ch
	}
	println("successNum = ", successNum)
	println("minPostTime = ", minPostTime)
	println("maxPostTime = ", maxPostTime)
	println("avgPostTime = ", totalTime/successNum)

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
	startTime := time.Now().UnixNano()

	var SERVERPOST string
	if len(SERVERPOSTS) == 0 {
		return
	} else if len(SERVERPOSTS) == 1 {
		SERVERPOST = SERVERPOSTS[0]
	} else {
		SERVERPOST = strconv.Itoa(rand.Intn(len(SERVERPOSTS) - 1))
	}
	body := common.DoPost(SERVERPOST+"/tx", data)

	if body == nil {
		postTx()
	}
	endTime := time.Now().UnixNano()
	postTime := endTime - startTime
	totalTime = + postTime
	if postTime < minPostTime {
		minPostTime = postTime
	}
	if postTime > maxPostTime {
		maxPostTime = postTime
	}
	successNum++
	//println(string(result))
	postTx()
}

func getPostTx() []byte {
	defer lock.Unlock()
	lock.Lock()
	i1 := l.Front()
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
