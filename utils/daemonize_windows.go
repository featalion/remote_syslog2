package utils

const CanDaemonize = false

func ResolvePath(path string) string {
	return path
}

func Daemonize(logFilePath, pidFilePath string, backTty *os.File) {
	panic("cannot daemonize on windows")
}
