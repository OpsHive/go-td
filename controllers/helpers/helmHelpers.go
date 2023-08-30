package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RetryHelmCommand(commandName string, args []string, kubeconfigPath string, maxRetries int, retryDelay time.Duration) ([]byte, error) {

	var output []byte
	var lastErr error

	for retry := 0; retry < maxRetries; retry++ {
		cmd := exec.Command(commandName, args...)
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)

		output, lastErr = cmd.CombinedOutput()

		if lastErr == nil {
			// Command succeeded, break out of the loop
			break
		}

		// Command failed, log the error and wait before retrying
		fmt.Printf("Attempt %d failed: %v\n", retry+1, lastErr)
		time.Sleep(retryDelay)
	}

	return output, lastErr
}
