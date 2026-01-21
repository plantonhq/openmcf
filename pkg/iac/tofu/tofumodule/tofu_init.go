package tofumodule

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tfbackend"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tfvars"
	"google.golang.org/protobuf/proto"
)

// Init initializes an HCL module (tofu or terraform) with optional JSON streaming.
// The binaryName parameter specifies which CLI binary to use ("tofu" or "terraform").
func Init(
	binaryName string,
	modulePath string,
	manifestObject proto.Message,
	backendType terraform.TerraformBackendType,
	backendConfigInput []string,
	providerConfigEnvVars []string,
	isJsonOutput bool,
	jsonLogEventsChan chan string,
) (err error) {
	if err := tfbackend.WriteBackendFile(modulePath, backendType); err != nil {
		return errors.Wrapf(err, "failed to write backend file")
	}

	tfVarsFile := filepath.Join(modulePath, ".terraform", "terraform.tfvars")
	if err := tfvars.WriteVarFile(manifestObject, tfVarsFile); err != nil {
		return errors.Wrapf(err, "failed to write %s file", tfVarsFile)
	}

	// Build the init command
	cmdArgs := []string{
		terraform.TerraformOperationType_init.String(),
		"--var-file", tfVarsFile,
	}
	if isJsonOutput {
		cmdArgs = append(cmdArgs, "-json")
	}
	for _, backendConfig := range backendConfigInput {
		cmdArgs = append(cmdArgs, "--backend-config", backendConfig)
	}

	cmd := exec.Command(binaryName, cmdArgs...)
	cmd.Dir = modulePath
	// https://stackoverflow.com/a/41133244
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, providerConfigEnvVars...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	fmt.Printf("%s module directory: %s\n", binaryName, modulePath)
	fmt.Printf("running command: %s\n", cmd.String())

	// If jsonLogEventsChan is provided, read stdout in a goroutine with panic recovery
	if jsonLogEventsChan != nil {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrap(err, "failed to create stdout pipe")
		}

		if err := cmd.Start(); err != nil {
			return errors.Wrapf(err, "failed to start %s command %s", binaryName, cmd.String())
		}

		errChan := make(chan error, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					panicErr := fmt.Errorf(
						"panic recovered in Init stdout reader: %v\nstack trace:\n%s",
						r, string(stack),
					)
					errChan <- panicErr
				}
				close(errChan)
			}()

			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				jsonLogEventsChan <- scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				errChan <- fmt.Errorf("error reading %s output: %v", binaryName, err)
			}
		}()

		if err := cmd.Wait(); err != nil {
			return errors.Wrapf(err, "failed to execute %s command %s", binaryName, cmd.String())
		}

		if readErr, ok := <-errChan; ok && readErr != nil {
			return readErr
		}
	} else {
		// Stream stdout to console
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "failed to execute %s command %s", binaryName, cmd.String())
		}
	}

	return nil
}
