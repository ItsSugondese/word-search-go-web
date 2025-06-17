package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func TailwindWatcher() {
	tailwindPath := "/usr/local/bin/tailwindcss" // from `which tailwindcss`

	cmd := exec.Command(tailwindPath,
		"-i", "./static/css/input.css",
		"-o", "./static/css/output.css",
		"--watch",
	)

	// Use your existing environment
	cmd.Env = os.Environ()

	// Pipe Tailwind output to your terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start tailwind watch: %v", err)
	}
	log.Println("Tailwind watch started...")

	// Optional: Handle CTRL+C or IDE Stop
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutting down Tailwind...")

		if err := cmd.Process.Kill(); err != nil {
			log.Printf("failed to kill tailwind process: %v", err)
		}
	}()

	// Let Go continue without blocking
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("tailwind exited: %v", err)
		}
	}()
}
