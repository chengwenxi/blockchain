package test

import (
	"testing"
	"os/exec"
	"fmt"
	"bytes"
	"strings"
	"time"
	"container/list"
	"sync"
	"strconv"
)

const (
	prefix = "test"
	pass   = "1234567890"
)

var ch = make(chan int)
var sendCh = make(chan int)
var l *list.List
var lock sync.Mutex

func Test_tx(t *testing.T) {

	//send to test account
	for i := 0; i < 5; i++ {
		go SendTx(i,5)
	}
	for i := 0; i < 5; i++ {
		<-sendCh
	}

	// resend to init1
	//println("start")
	//l = AccountList()
	//println(l.Len())
	//for i := 0; i < 300; i++ {
	//	go resend()
	//}
	//for i := 0; i < 300; i++ {
	//	<-ch
	//}
	//println("end")
}

func SendTx(i int,gos int) {
	kv := AccountMap()
	number := int(5000/gos)
	for j := 0; j < number; j++ {
		c := exec.Command("iris", "client", "tx", "send", "--to="+kv["test"+strconv.Itoa(3000+number*i+j)],
			"--amount=10000iris", "--name="+"cwx"+strconv.Itoa(i+1), "--password="+pass)
		var out bytes.Buffer
		c.Stdout = &out
		println(strconv.Itoa(3000+number*i+j))
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("%s", out.String())
		time.Sleep(5 * time.Second)
	}
	sendCh <- 0
}

func resend() {
	name := getName()
	if name != "" {
		c := exec.Command("iris", "client", "tx", "send", "--to=CAF62CF4258BB500D91C775106AD6419986B2A94",
			"--amount=1iris", "--name="+name, "--password="+pass)
		var out bytes.Buffer
		c.Stdout = &out
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		//fmt.Printf("%s", out.String())
		resend()
		//time.Sleep(1 * time.Second)
	}
	ch <- 0
}

func getName() string {
	defer lock.Unlock()
	lock.Lock()
	//i1 := l.Front()
	i1 := l.Back()
	name := ""
	if i1 != nil {
		name = i1.Value.(string)
		l.Remove(i1)
	}
	return name
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
			if len(l) > 1 && strings.Index(l[0], prefix) != -1 {
				kv[l[0]] = l[1]
			}
		}
	}
	return kv
}

func AccountList() *list.List {
	c := exec.Command("iris", "client", "keys", "list")
	var out bytes.Buffer
	c.Stdout = &out
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	result := out.String()
	s := strings.Split(result, "\n")
	accounts := list.New()
	for i, v := range s {
		if i == 0 {
			continue
		} else {
			l := strings.Split(v, "\t\t")
			if len(l) > 1 && strings.Index(l[0], prefix) != -1 {
				accounts.PushBack(l[0])
			}
		}
	}
	return accounts
}
