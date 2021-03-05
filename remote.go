package utility

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"time"
)

func SecureCopyLinux(src string, dest string, timeout_sec int) error {
	if runtime.GOOS != "linux" {
		return errors.New("securecopylinux: is only supported on unix systems")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout_sec)*time.Second)
	defer cancel()
	_, err := exec.CommandContext(ctx, "scp", "-o ConnectTimeout=30", src, dest).Output()
	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}
	if err != nil {
		return err
	}
	return nil
}
