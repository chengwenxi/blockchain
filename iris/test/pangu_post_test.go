package post

import (
	"testing"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"net/http"
	"log"
	"container/list"
	"sync"
	"time"
)

var TO = "CAF62CF4258BB500D91C775106AD6419986B2A94"
var PASSWORD = "1234567890"
var SERVER = "http://116.62.62.39:1198"
var ch = make(chan int)
var signCh = make(chan int)
var l = list.New()
var lock sync.Mutex

var goNum = 200


func Test_PostTx(t *testing.T) {
	keys := getKeys()
	for i := 0; i < goNum; i++ {
		j := len(keys)/goNum
		go buildAndSignTxAll(keys, j*i, j)
	}
	for i := 0; i < goNum; i++ {
		<- signCh
	}
	println("sign end, number = " ,l.Len())
	startTime := time.Now().Unix()
	for i := 0; i < goNum; i++ {
		go postTx()
	}
	for i := 0; i < goNum; i++ {
		<- ch
	}
	endTime := time.Now().Unix()
	println(endTime-startTime)
}

func buildAndSignTxAll(keys []Key, start int, number int) {
	for i := start; i < start+number; i++ {
		buildAndSignTx(keys[i].Name, keys[i].Address)
	}
	signCh <- 0
}

func postTx()  {
	data := getPostTx()
	if data==nil{
		ch <-0
		return
	}
	_ = DoPost(SERVER+"/tx", data)
	//println(string(result))
	postTx()
}

func getPostTx() []byte{
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

func getKeys() []Key {
	result := DoGet(SERVER + "/keys")
	var keys []Key
	if result != nil {
		json.Unmarshal(result, &keys)
	}
	return keys
}

func buildAndSignTx(name string, addr string) {

	amount := 1
	coin := "iris"
	result := DoGet(SERVER + "/query/nonce/" + addr)
	var nonce Nonce
	if result != nil {
		json.Unmarshal(result, &nonce)
	}
	nonce.Data += 1

	//build send
	si := new(SendInput)
	si.Amount = Coins{Coin{Denom: coin, Amount: int64(amount)}}
	si.From = &Actor{ChainID: "", App: "sigs", Address: addr}
	si.To = &Actor{ChainID: "", App: "sigs", Address: TO}
	//si.Fees = &Coin{Denom: feeCoin, Amount: 1}
	si.Sequence = nonce.Data
	siStr, _ := json.Marshal(si)
	result = DoPost(SERVER+"/build/send", siStr)
	if result == nil {
		return
	}

	//sign tx
	requestSign := new(RequestSign)
	requestSign.Name = name
	requestSign.Password = PASSWORD
	json.Unmarshal(result, &requestSign.Tx)
	rsStr, _ := json.Marshal(requestSign)
	result = DoPost(SERVER+"/sign", rsStr)
	if result == nil {
		return
	}
	//println(string(result))
	lock.Lock()
	l.PushBack(result)
	lock.Unlock()
	//send tx
}

type Key struct {
	Name    string
	Address string
}

type Nonce struct {
	Height uint32
	Data   uint32
}

type Actor struct {
	ChainID string `json:"chain"` // this is empty unless it comes from a different chain
	App     string `json:"app"`   // the app that the actor belongs to
	Address string `json:"addr"`  // arbitrary app-specific unique id
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

type Coins []Coin

type RequestSign struct {
	Name     string `json:"name,omitempty" validate:"required,min=3,printascii"`
	Password string `json:"password,omitempty" validate:"required,min=10"`

	Tx map[string]interface{} `json:"tx" validate:"required"`
}

type SendInput struct {
	Fees     *Coin  `json:"fees"`
	Multi    bool   `json:"multi,omitempty"`
	Sequence uint32 `json:"sequence"`

	To     *Actor `json:"to"`
	From   *Actor `json:"from"`
	Amount Coins  `json:"amount"`
}

type Result struct {
	Code  uint32
	Error string
}

func DoGet(url string) []byte {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	var result Result
	json.Unmarshal(body, &result)
	if result.Error != "" {
		log.Println(result.Error)
		return nil
	}
	return body
}

func DoPost(url string, data []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	var result Result
	json.Unmarshal(body, &result)
	if result.Error != "" {
		log.Println(result.Error)
		return nil
	}
	return body
}
