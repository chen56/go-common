package iox

import (
	"os"
	"bufio"
)

//ref: http://learngowith.me/a-better-way-of-handling-stdin/
func IsPipeMode() bool {
	// Retrieve file information of stdin (os.FileInfo)
	stdinFileInfo, _ := os.Stdin.Stat()

	// Print string representation of stdin's file mode.
	//fmt.Println(stdinFileInfo.Mode().String())

	if (stdinFileInfo.Mode()&os.ModeNamedPipe != 0) {
		return true
	}
	return false
}

func ReadArgsIfPipe(args []string) []string {
	if (IsPipeMode()) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			str := scanner.Text()
			args = append(args, str)
		}
	}
	return args
}
