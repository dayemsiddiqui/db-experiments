package docker

import (
	"log"
	"os/exec"
	"strings"
)

func IsContainerRunning() bool {
	cmd := exec.Command("docker-compose", "ps", "-q", "postgres")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

func RunContainer() {
	if IsContainerRunning() {
		cmd := exec.Command("docker-compose", "up", "-d")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to start the Postgres container: %s", err)
		}
	}
}
