package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("curl", "http://18.162.47.230/queue/api/json", "-u", "admin:lihongzhi", "|", "grep", "work-node")
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}
