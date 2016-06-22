package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func daemon() {
	//设置标准输入输出
	files := make([]*os.File, 3)
	files[0] = os.Stdin
	files[1] = os.Stdout
	if *stdout != "" {
		files[1], _ = os.OpenFile(*stderr, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	}
	files[2] = os.Stderr
	if *stderr != "" {
		files[2], _ = os.OpenFile(*stdout, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	}

	//设置工作目录
	var dir string
	if *wdir != "" {
		dir = *wdir
	} else {
		dir, _ = os.Getwd()
	}

	//设置程序执行用户
	var credential *syscall.Credential
	if *guser != "" {
		u, err := user.Lookup(*guser)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if u.Uid != "" {
			uid, _ := strconv.Atoi(u.Uid)
			gid, _ := strconv.Atoi(u.Gid)
			credential = &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			}
		}
	}

	//设置程序环境变量
	var envs = os.Environ()
	if *environment != "" {
		for _, e := range strings.Split(*environment, ",") {
			envs = append(envs, strings.TrimSpace(e))
		}
	}

	sysattrs := syscall.SysProcAttr{
		//设置使用session
		Setsid:     true,
		Credential: credential,
	}

	attrs := os.ProcAttr{
		Dir:   dir,
		Env:   envs,
		Files: files, Sys: &sysattrs,
	}

	//最后一次fork
	proc, err := os.StartProcess(*exec, os.Args[inx-1:], &attrs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create process %s", err)
		os.Exit(2)
	}
	proc.Release()

	if *pidfile != "" {
		err := ioutil.WriteFile(*pidfile, []byte(strconv.Itoa(proc.Pid)), 0660)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't write pidfile %s: %s", *pidfile, err)
			os.Exit(2)
		}
	}

}
