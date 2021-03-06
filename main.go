package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"io/ioutil"
	"strconv"
)

//docker         run image <cmd> <params>
//go run main.go run       <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Println("runnning", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run();
}

func child() {
	fmt.Println("runnning", os.Args[2:])
	
	cg()

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/home/ubuntufs")
	syscall.Chdir("/")
	syscall.Mount("proc","proc","proc",0,"")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run();

	syscall.Unmount("/proc",0)
}

func cg(){
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	must(os.Mkdir(filepath.Join(pids, "yash"), 0755))
	must(ioutil.WriteFile(filepath.Join(pids, "yash/pids.max"), []byte("20"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "yash/notify_on_release"), []byte("1"),0700))
	must(ioutil.WriteFile(filepath.Join(pids, "yash/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

}

func must(err error){
	if err != nil{
		panic(err)
	}
}
