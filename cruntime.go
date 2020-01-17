package main

import (
	"fmt"
	"io/ioutil"
	"linux-container/utils"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var (
	FORK_ROUTINE             = []string{"child"}
	CONTAINER_NAME           = "debian-fs"
	CONTAINER_ROOT_FS_PATH   = "./debian-rootfs/build/amd64/amd64-rootfs-20200114T195115Z"
	CONTAINER_ROOT_PATH      = "/"
	CONTAINER_PROC_PATH      = "proc"
	CGROUP_MAX_PID           = "pids.max"
	CGROUP_NOTIFY_ON_RELEASE = "notify_on_release"
	CGROUP_PROCS             = "cgroup.procs"
	CGROUP_FS_PATH           = "/sys/fs/cgroup"
	CGROUP_PIDS              = "pids"
	NETSETGO_PATH = "/usr/local/bin/netsetgo"
)

func main() {
	if len(os.Args) < 2 {
		panic(utils.Red("No command!"))
	}
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic(utils.Red("Unknown command!"))
	}
}

func run() {
	WelcomeSession()
	ForkWithNSFlags()
}

func child() {
	EnforceCGroups()
	SetNewUTS()
	ChrootAndChpath()
	SetNewNS()
	WelcomeSession()
	ContainerProcess()
	UnsetNewNS()
	WaitForNetwork()
}

func WelcomeSession() {
	fmt.Printf(utils.Magenta("Running %v as pid:%d, uid:%d and gid:%d\n"), os.Args[2:], os.Getpid(), os.Getuid(), os.Getgid())
}

func ForkWithNSFlags() {
	exe, _ := os.Executable()
	cmd := exec.Command(exe, append(FORK_ROUTINE, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// set flags for the 6 namespace restrictions
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	ensure(cmd.Start())
	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	netSetGoCmd := exec.Command(NETSETGO_PATH, "-pid", pid)
	if err := netSetGoCmd.Run(); err != nil {
		fmt.Printf("Error running netsetgo - %s\n", err)
		os.Exit(1)
	}
	ensure(cmd.Wait())
}

func WaitForNetwork() error {
	maxWait := time.Second * 3
	checkInterval := time.Second
	timeStarted := time.Now()

	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}

		// pretty basic check ...
		// > 1 as a lo device will already exist
		if len(interfaces) > 1 {
			return nil
		}

		if time.Since(timeStarted) > maxWait {
			return fmt.Errorf("Timeout after %s waiting for network", maxWait)
		}

		time.Sleep(checkInterval)
	}
}

func EnforceCGroups() {
	// Set max container processes to 20
	pids := filepath.Join(CGROUP_FS_PATH, CGROUP_PIDS)
	ensure(os.Mkdir(filepath.Join(pids, CONTAINER_NAME), 0755))
	ensure(ioutil.WriteFile(filepath.Join(pids, CONTAINER_NAME, CGROUP_MAX_PID), []byte("20"), 0700))
	ensure(ioutil.WriteFile(filepath.Join(pids, CONTAINER_NAME, CGROUP_NOTIFY_ON_RELEASE), []byte("1"), 0700))
	ensure(ioutil.WriteFile(filepath.Join(pids, CONTAINER_NAME, CGROUP_PROCS), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func SetNewUTS() {
	ensure(syscall.Sethostname([]byte(CONTAINER_NAME)))
}

func ChrootAndChpath() {
	ensure(syscall.Chroot(CONTAINER_ROOT_FS_PATH))
	ensure(syscall.Chdir(CONTAINER_ROOT_PATH))
}

func UnsetNewNS() {
	ensure(syscall.Unmount(CONTAINER_PROC_PATH, 0))
}

func SetNewNS() {
	ensure(syscall.Mount(CONTAINER_PROC_PATH, CONTAINER_PROC_PATH, CONTAINER_PROC_PATH, 0, ""))
}

func ContainerProcess() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ensure(cmd.Run())
}

func ensure(err error) {
	if err != nil && !os.IsExist(err) {
		panic(utils.Red(err))
	}
}
