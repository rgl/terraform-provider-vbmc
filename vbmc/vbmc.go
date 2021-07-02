package vbmc

import (
	"bytes"
	"encoding/json"
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

func execVbmc(args ...string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command("vbmc", args...)
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
	_, err = execVbmc(
		"add",
		domainName,
		"--address", address,
		"--port", strconv.Itoa(port),
		"--username", username,
		"--password", password)
	if err != nil {
		return nil, err
	}
	_, err = execVbmc("start", domainName)
	if err != nil {
		return nil, err
	}
	return Get(domainName)
}

func Delete(domainName string) error {
	_, err := execVbmc("delete", domainName)
	if err != nil {
		if execError, ok := err.(*VbmcExecError); ok {
			if strings.Contains(execError.Stderr, "No domain with matching name") {
				return nil
			}
		}
		return err
	}
	return nil
}

func Get(domainName string) (*Vbmc, error) {
	data, err := execVbmc("show", "-f", "json", "--noindent", domainName)
	if err != nil {
		return nil, err
	}

	var properties []struct {
		Property string
		Value    json.RawMessage // NB the Value type depends on the Property. It can be a string, number, etc.
	}

	if err := json.Unmarshal([]byte(data), &properties); err != nil {
		return nil, fmt.Errorf("failed to parse vbmc output: %v", err)
	}

	vbmc := &Vbmc{}

	for _, property := range properties {
		switch property.Property {
		case "domain_name":
			if err := json.Unmarshal(property.Value, &vbmc.DomainName); err != nil {
				return nil, fmt.Errorf("failed to parse vbmc output domain_name: %v", err)
			}
		case "port":
			if err := json.Unmarshal(property.Value, &vbmc.Port); err != nil {
				return nil, fmt.Errorf("failed to parse vbmc output port: %v", err)
			}
		}
	}

	return vbmc, nil
}
