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
)

const (
	prefix = "test"
	pass   = "1234567890"
	pass1  = "1111111111"
)

var ch = make(chan int)
var l *list.List
var lock sync.Mutex

func Test_tx(t *testing.T) {

	//SendTx()
	println("start")
	l = AccountList()
	println(l.Len())
	for i := 0; i < 300; i++ {
		go resend()
	}
	for i := 0; i < 300; i++ {
		<-ch
	}
	println("end")
}

func SendTx() {
	kv := AccountMap()
	fmt.Print(len(kv))
	for _, v := range kv {
		c := exec.Command("iris", "client", "tx", "send", "--to="+v,
			"--amount=10000iris", "--name=init1", "--password="+pass1)
		var out bytes.Buffer
		c.Stdout = &out
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("%s", out.String())
		time.Sleep(5 * time.Second)
	}
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
