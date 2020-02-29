package ffmpeg

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"time"
)

func runWithTimeout(command string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return runContext(ctx, command)
}

func runContext(ctx context.Context, command string) ([]byte, error) {
	var cmd *exec.Cmd
	var err error
	command = binPath + " " + command

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf

	err = cmd.Start()
	if err == exec.ErrNotFound {
		return []byte{}, ErrBinNotFound
	} else if err != nil {
		return []byte{}, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		err = cmd.Process.Kill()
		if err == nil {
			return []byte{}, ErrTimeout
		}
		return []byte{}, err
	case err = <-done:
		if err != nil {
			return []byte{}, err
		}
	}

	if err != nil {
		return outputBuf.Bytes(), err
	}

	return outputBuf.Bytes(), nil
}
