package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("java", "-jar", "jenkins-cli.jar", "-s", "http://18.163.184.77:80", "delete-node", "test")
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}
