// +build !windows

package systemcall

import (
	"log"
	"syscall"
)

func SetFileLimit() {
	var err error
	var rlim syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		log.Fatal("get rlimit error: " + err.Error())
		return
	}
	rlim.Cur = 10240
	rlim.Max = 10240
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		log.Fatal("set rlimit error: " + err.Error())
		return
	}
}
