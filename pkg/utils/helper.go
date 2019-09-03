package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetHostName get host name
func GetHostName() string {
	hostname, _ := RunCommandOnHost("cat", "/etc/hostname")
	return strings.TrimSuffix(string(hostname), "\n")
}

// GetAPIServerFQDN gets the API Server FQDN from the kubeconfig file
func GetAPIServerFQDN() (string, error) {
	output, err := RunCommandOnHost("cat", "/var/lib/kubelet/kubeconfig")

	if err != nil {
		log.Println("Can't open kubeconfig file: ", err)
		return "", err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		index := strings.Index(line, "server: ")
		if index >= 0 {
			fqdn := line[index+len("server: "):]
			fqdn = strings.Replace(fqdn, "https://", "", -1)
			fqdn = strings.Replace(fqdn, ":443", "", -1)
			return fqdn, nil
		}
	}

	return "", errors.New("Could not find server definitions in kubeconfig")
}

// GetAzureBlobCredential get azure blob access info
func GetAzureBlobCredential() (string, string) {
	accountName, _ := ioutil.ReadFile("/etc/azure-blob/accountName")
	sasKey, _ := ioutil.ReadFile("/etc/azure-blob/sasKey")
	return strings.TrimSuffix(string(accountName), "\n"), strings.TrimSuffix(string(sasKey), "\n")
}

// RunCommandOnHost runs a command on host system
func RunCommandOnHost(command string, arg ...string) (string, error) {
	args := []string{"--target", "1", "--mount", "--uts", "--ipc", "--net", "--pid"}
	args = append(args, "--")
	args = append(args, command)
	args = append(args, arg...)

	cmd := exec.Command("nsenter", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// RunCommandOnContainer runs a command on container system
func RunCommandOnContainer(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// WriteToFile writes data to a file
func WriteToFile(fileName string, data string) error {
	f, _ := os.Create(fileName)
	defer f.Close()

	_, err := f.Write([]byte(data))

	return err
}

// CreateCollectorDir creates a working dir for a collector
func CreateCollectorDir(name string) (string, error) {
	rootPath := filepath.Join("/aks-diagnostic/", GetHostName(), "metrics", name)
	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return rootPath, nil
}

// CreateDiagnosticDir creates a working dir for diagnostic
func CreateDiagnosticDir() (string, error) {
	rootPath := filepath.Join("/aks-diagnostic/", GetHostName(), "diagnostic")
	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return rootPath, nil
}
