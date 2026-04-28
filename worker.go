package main

import (
	"bufio"
	"fmt"
	"image/png"
	"net"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/vova616/screenshot"
)

func setWallpaper(path string) {
	uintPath, _ := syscall.UTF16PtrFromString(path)
	uintAction := uintptr(20)         // SPI_SETDESKWALLPAPER
	uintParam := uintptr(0x01 | 0x02) // SPIF_UPDATEINIFILE | SPIF_SENDWININICHANGE

	user32 := syscall.NewLazyDLL("user32.dll")
	systemParametersInfo := user32.NewProc("SystemParametersInfoW")
	systemParametersInfo.Call(uintAction, 0, uintptr(unsafe.Pointer(uintPath)), uintParam)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	fullCommand := strings.TrimSpace(message)
	fmt.Printf("Received: [%s]\n", fullCommand)

	switch {
	case fullCommand == "lock":
		exec.Command("rundll32.exe", "user32.dll,LockWorkStation").Run()
	case fullCommand == "shutdown":
		exec.Command("shutdown", "/s", "/t", "60").Run()
	case strings.HasPrefix(fullCommand, "wallpaper "):
		imagePath := strings.TrimPrefix(fullCommand, "wallpaper ")
		setWallpaper(imagePath)
	case fullCommand == "screenshot":
		img, err := screenshot.CaptureScreen()
		if err == nil {
			// Send the image back through the connection
			png.Encode(conn, img)
		}
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Agent is listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}