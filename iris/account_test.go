package test

import (
	"testing"
	"os/exec"
	"strconv"
	"bytes"
	"fmt"
)

func Test_NewAccount(t *testing.T) {
	for i := 0; i < number; i++ {
		c := exec.Command("iris", "client", "keys", "new", prefix+strconv.Itoa(i), "--password="+pass)
		var out bytes.Buffer
		c.Stdout = &out
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("%s", out.String())
	}
}
