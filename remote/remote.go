package remote

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// SecureCopyLinux копирует файл по сети с использованием scp. Работает только в системе Unix. Таймаут указывается
// в секундах. Также нужно передать указатели для блокировки ресурсов при совместном доступе или nil при отсутствии
// необходимости блокировки.
func SecureCopyLinux(src string, dest string, timeout int, srcLocker *sync.RWMutex, destLocker *sync.RWMutex) error {
	if runtime.GOOS != "linux" {
		return errors.New("SecureCopyLinux: is only supported on unix systems")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if srcLocker != nil {
		srcLocker.RLock()
		defer srcLocker.RUnlock()
	}
	if destLocker != nil {
		destLocker.Lock()
		defer destLocker.Unlock()
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
