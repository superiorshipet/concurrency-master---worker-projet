package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func sendCommand(ip string, command string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", ip, 5*time.Second)
	if err != nil {
		fmt.Printf("[-] Error connecting to %s: %v\n", ip, err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, command+"\n")

	if command == "screenshot" {
		fileName := fmt.Sprintf("screenshot_%s.png", strings.ReplaceAll(ip, ":8080", ""))
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("[-] Failed to create file for %s\n", ip)
			return
		}
		defer file.Close()
		_, err = io.Copy(file, conn)
		if err != nil {
			fmt.Printf("[-] Failed to save screenshot from %s\n", ip)
		} else {
			fmt.Printf("[+] Screenshot saved: %s\n", fileName)
		}
	} else {
		fmt.Printf("[+] Command '%s' sent to %s\n", command, ip)
	}
}

func main() {
	agents := []string{
		"10.42.0.41:8080",
		"10.42.0.182:8080",
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Distributed Controller ---")
		fmt.Println("Commands: lock | shutdown | wallpaper | screenshot | exit")
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		if command == "exit" {
			break
		}

		if command == "wallpaper" {
			fmt.Print("Enter Agent Image Path: ")
			path, _ := reader.ReadString('\n')
			command = "wallpaper " + strings.TrimSpace(path)
		}

		var wg sync.WaitGroup
		for _, ip := range agents {
			wg.Add(1)
			go sendCommand(ip, command, &wg)
		}
		wg.Wait()
	}
}
