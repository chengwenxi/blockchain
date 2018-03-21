package test

import (
	"testing"
	"os/exec"
	"fmt"
	"bytes"
	"strings"
)

const (
	number = 100
	prefix = "test"
	pass   = "1234567890"
	pass1  = "1111111111"
)

var ch = make(chan int)

func Test_tx(t *testing.T) {

	kv := AccountList()
	fmt.Print(len(kv))
	//for k, v := range kv {
	//	go SendTx(k, v)
	//}
	//for i := 0; i < len(kv); i++ {
	//	<-ch
	//}
}

func SendTx(name string, addr string) {
	c := exec.Command("iris", "client", "tx", "send", "--to "+addr,
		"--amount 100iris", "--name init1", "--password="+pass1)
	var out bytes.Buffer
	c.Stdout = &out
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Printf("%s", out.String())
	for i := 0; i < 100; i++ {
		c1 := exec.Command("iris", "client", "tx", "send", "--to CAF62CF4258BB500D91C775106AD6419986B2A94",
			"--amount 1iris", "--name"+name, "--password="+pass)
		var out1 bytes.Buffer
		c1.Stdout = &out1
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("%s", out1.String())
	}
	ch <- 0
}

func AccountList() map[string]string {
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
			if len(l) > 1 && strings.Index(l[0], "test") != -1 {
				kv[l[0]] = l[1]
			}
		}
	}
	return kv
}
