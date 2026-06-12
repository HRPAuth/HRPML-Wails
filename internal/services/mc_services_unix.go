//go:build !windows

package services

import (
	"os/exec"
	"syscall"
)

// configureDetachedProcess puts the child in its own process group (Setpgid)
// so it survives when the parent (Wails app) exits. On Linux this also makes
// it possible to later signal the whole game tree via kill(-pid, ...).
func configureDetachedProcess(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}
