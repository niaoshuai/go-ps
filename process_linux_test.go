// +build linux

package ps

import (
	"testing"
)

// 获取进程列表
func TestLinuxProcesses(t *testing.T) {
	processesList, _ := processes()

	for _, item := range processesList {
		t.Log(item)
	}
}

func TestLinuxCmdLine(t *testing.T) {
	uni := UnixProcess{}
	uni.pid = 9104
	err := uni.CmdLine()
	if err != nil {
		t.Error(err)
	}
	t.Log(uni.cmd)
	t.Log(uni.env)
}
