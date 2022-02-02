package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func macroListener() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("[SCIMPLISTIC] Command: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		stop := handleMacro(text)
		if stop {
			return
		}
	}
}

func handleMacro(macro string) bool {
	if macro == "pause" {
		Pause()
	} else if macro == "resume" {
		Resume()
	} else if macro == "shutdown" {
		ShutDown()
		return true
	} else if macro == "version" {
		fmt.Println(version)
	}
	return false
}
