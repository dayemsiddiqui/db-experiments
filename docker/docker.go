package docker

import (
	"log"
	"os/exec"
)

func IsContainerNotRunning() bool {
	log.Println("Checking if Postgres container is running...")
	cmd := exec.Command("docker", "compose", "ps", "-q", "postgres")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to check if Postgres container is running: %s", err)
		return false
	}
	return len(output) == 0
}

func RunContainer() {
	if IsContainerNotRunning() {
		log.Println("Starting Postgres container...")
		cmd := exec.Command("docker", "compose", "up", "-d")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to start the Postgres container: %s", err)
		}
	} else {
		log.Println("Postgres container is already running.")
	}
}
