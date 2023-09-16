package service

import (
	"fmt"
	"net/rpc"
	"testing"
)

func TestMaster_ServerTest(t *testing.T) {
	c, err := rpc.Dial("tcp", "localhost:8848")
	if err != nil {
		fmt.Println(" rpc.Dial", err)
		// Master结束进程，退出worker
	}
	err = c.Call("Master.AssignTask", "", "")
	if err != nil {
		fmt.Println("err", err)
	}
}
