package browser

import (
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
)

func OpenURL(url string) error {
	var cmd string
	var args []string
	slog.Debug("opening default browser to %s", slog.String("url", url))

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func PrintableOpenURL(url string) string {
	return fmt.Sprintf("# Open this URL in the browser: %s", url)
}
