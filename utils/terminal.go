package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// ClearTerminal cleares the current terminal window.
func ClearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("Error while clearing the terminal")
	}
}
