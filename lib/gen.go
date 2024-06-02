//go:generate go run gen.go

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if err := darwinDownload(); err != nil {
		log.Fatalf("failed to download file for macos: %v", err)
	}
	if err := linuxDownload(); err != nil {
		log.Fatalf("failed to download file for linux: %v", err)
	}
}

func darwinDownload() error {
	const path = "/opt/homebrew/lib"
	const name = "libemulator.dylib"

	log.Println("starting download lib for macos")

	initTonCmd := exec.Command("brew", "tap", "ton-blockchain/ton")
	if output, err := initTonCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("[darwinDownload] failed to init ton: %v, output: %s", err, output)
	}

	_, err := os.Stat(fmt.Sprintf("%v/%v", path, name))
	if err == nil {
		log.Println("file already exist, reinstalling lib...")
		reinstallCmd := exec.Command("brew", "reinstall", "ton")
		if output, err := reinstallCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("[darwinDownload] failed to reinstall lib: %v, output: %s", err, output)
		}
	} else if os.IsNotExist(err) {
		log.Println("file doesn't exist, installing lib...")
		installCmd := exec.Command("brew", "install", "ton")
		if output, err := installCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("[darwinDownload] failed to install lib: %v, output: %s", err, output)
		}
	} else {
		return fmt.Errorf("failed to check file: %v", err)
	}

	copyCmd := exec.Command("cp", fmt.Sprintf("%v/%v", path, name), "darwin/")
	if output, err := copyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("[darwinDownload] failed to copy file: %v, output: %s", err, output)
	}

	log.Println("[darwinDownload] successfully update lib")

	return nil
}

func linuxDownload() error {
	const path = "/usr/lib"
	const name = "libemulator.so"

	log.Println("starting download lib for linux")

	commands := [][]string{
		{"sudo", "apt-key", "adv", "--keyserver", "keyserver.ubuntu.com", "--recv-keys", "F6A649124520E5F3"},
		{"sudo", "add-apt-repository", "ppa:ton-foundation/ppa"},
		{"sudo", "apt", "update"},
		{"sudo", "apt", "install", "ton"},
	}
	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("[linuxDownload] failed to install lib: %v, output: %s", err, output)
		}
	}

	copyCmd := exec.Command("cp", fmt.Sprintf("%v/%v", path, name), "linux/")
	if output, err := copyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("[linuxDownload] failed to copy file: %v, output: %s", err, output)
	}

	log.Println("[linuxDownload] successfully update lib")

	return nil
}
