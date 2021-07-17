package vbmc

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type VbmcExecError struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func (err *VbmcExecError) Error() string {
	return fmt.Sprintf("failed to exec vbmc: exitCode=%d stdout=%s stderr=%s", err.ExitCode, err.Stdout, err.Stderr)
}

type Vbmc struct {
	DomainName string
	Port       int
}

func getContainerName(domainName string) string {
	return fmt.Sprintf("vbmc-emulator-%s", domainName)
}

func docker(args ...string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command("docker", args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		exitCode := -1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ProcessState.ExitCode()
		}
		return "", &VbmcExecError{
			ExitCode: exitCode,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
		}
	}

	return strings.TrimSpace(stdout.String()), nil
}

func Create(domainName string, address string, port int, username string, password string) (*Vbmc, error) {
	err := Delete(domainName)
	if err != nil {
		return nil, err
	}
	_, err = docker(
		"run",
		"--rm",
		"--name",
		getContainerName(domainName),
		"--detach",
		"-v",
		"/var/run/libvirt/libvirt-sock:/var/run/libvirt/libvirt-sock",
		"-v",
		"/var/run/libvirt/libvirt-sock-ro:/var/run/libvirt/libvirt-sock-ro",
		"-e",
		fmt.Sprintf("VBMC_EMULATOR_DOMAIN_NAME=%s", domainName),
		"-e",
		fmt.Sprintf("VBMC_EMULATOR_USERNAME=%s", username),
		"-e",
		fmt.Sprintf("VBMC_EMULATOR_PASSWORD=%s", password),
		"-p",
		fmt.Sprintf("%s:%d:6230", address, port),
		"ruilopes/vbmc-emulator")
	if err != nil {
		return nil, err
	}
	vbmc, err := Get(domainName)
	if err != nil {
		return nil, err
	}
	if vbmc == nil {
		return nil, fmt.Errorf("failed to create the vbmc container; it probably died for unknown reasons")
	}
	return vbmc, nil
}

func Delete(domainName string) error {
	containerName := getContainerName(domainName)
	_, err := docker("kill", "--signal", "INT", containerName)
	if err != nil {
		if execError, ok := err.(*VbmcExecError); ok {
			if strings.Contains(execError.Stderr, "No such container") {
				return nil
			}
		}
		return err
	}
	_, err = docker("wait", containerName)
	if err != nil {
		if execError, ok := err.(*VbmcExecError); ok {
			if strings.Contains(execError.Stderr, "No such container") {
				return nil
			}
		}
		return err
	}
	return nil
}

func Get(domainName string) (*Vbmc, error) {
	stdout, err := docker("port", getContainerName(domainName), "6230")
	if err != nil {
		if execError, ok := err.(*VbmcExecError); ok {
			if strings.Contains(execError.Stderr, "No such container") {
				return nil, nil
			}
		}
		return nil, err
	}

	vbmc := &Vbmc{
		DomainName: domainName,
	}

	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		// e.g. 0.0.0.0:6230
		line := scanner.Text()
		parts := strings.SplitN(line, ":", -1)
		if len(parts) < 1 {
			continue
		}
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			return nil, err
		}
		vbmc.Port = port
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return vbmc, nil
}
