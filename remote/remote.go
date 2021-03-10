package remote

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func SecureCopyLinux(src string, dest string, timeout_sec int, src_locker *sync.RWMutex, dest_locker *sync.RWMutex) error {
	if runtime.GOOS != "linux" {
		return errors.New("SecureCopyLinux: is only supported on unix systems")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout_sec)*time.Second)
	defer cancel()
	if src_locker != nil {
		src_locker.RLock()
		defer src_locker.RUnlock()
	}
	if dest_locker != nil {
		dest_locker.Lock()
		defer dest_locker.Unlock()
	}
	_, err := exec.CommandContext(ctx, "scp", "-o ConnectTimeout=30", src, dest).Output()
	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}
	if err != nil {
		return err
	}
	return nil
}
