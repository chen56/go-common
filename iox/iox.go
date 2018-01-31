package iox

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

//ref: http://learngowith.me/a-better-way-of-handling-stdin/
func IsPipeMode() bool {
	// Retrieve file information of stdin (os.FileInfo)
	stdinFileInfo, _ := os.Stdin.Stat()
	// Print string representation of stdin's file mode.
	//fmt.Println(stdinFileInfo.Mode().String())
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		return true
	}
	return false
}

func ReadFromPipe() string {
	if !IsPipeMode() {
		return ""
	}
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
	}
	return buf.String()
}

func SafeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
