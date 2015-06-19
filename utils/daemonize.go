// +build !windows

package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/VividCortex/godaemon"
	"github.com/leonsodhi/lockfile"
)

const CanDaemonize = true

func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(os.Getenv("__DAEMON_CWD"), path)
}

func Daemonize(logFilePath, pidFilePath string, backTty *os.File) {
	if os.Getenv("__DAEMON_CWD") == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot determine working directory: %v", err)
			os.Exit(1)
		}
		os.Setenv("__DAEMON_CWD", cwd)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not open local log file: %v", err)
		os.Exit(1)
	}

	stdout, stderr, err := godaemon.MakeDaemon(&godaemon.DaemonAttr{CaptureOutput: true, Files: []**os.File{&backTty}})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not Daemonize: %v", err)
		os.Exit(1)
	}

	proxyOutputFn := func(rd io.Reader) {
		if backTty != nil {
			io.Copy(backTty, rd)
		} else {
			io.Copy(logFile, rd)
		}
	}
	go proxyOutputFn(stdout)
	go proxyOutputFn(stderr)

	lock, err := lockfile.New(pidFilePath)
	err = lock.TryLock()
	if err != nil {
		fmt.Println("Cannot lock \"%v\": %v", lock, err)
		os.Exit(1)
	}

}
