package genfilops

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// handlErr is a shitty way to deal with errors.
func handlErr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func stringToFileMode(inString string) os.FileMode {
	permUint, err := strconv.ParseUint(inString, 8, 64)
	handlErr(err)
	return os.FileMode(permUint)
}

func ReadFileToArray(fileName string) []string {
	var retSlice []string

	srcFile, err := os.Open(fileName)
	handlErr(err)
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)

	for scanner.Scan() {
		retSlice = append(retSlice, scanner.Text())
	}

	handlErr(scanner.Err())

	return retSlice
}

func WriteBytesToFile(fileName string, outputBytes []byte, filePerm string) int {
	permFM := stringToFileMode(filePerm)
	outFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, permFM)
	handlErr(err)
	defer outFile.Close()

	writtenBytes, err := outFile.Write(outputBytes)
	handlErr(err)

	return writtenBytes
}

func WriteStringsToFile(fileName string, outputStrings []string, filePerm string) int {
	permFM := stringToFileMode(filePerm)
	outFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, permFM)
	handlErr(err)
	defer outFile.Close()

	writeString := strings.Join(outputStrings, "\n")
	writtenBytes, err := outFile.WriteString(writeString)
	handlErr(err)

	return writtenBytes
}

func SetFileTimestamps(fileName string, accessTime string, modTime string) {
	aTime, err := time.Parse(time.RFC3339, accessTime)
	handlErr(err)
	mTime, err := time.Parse(time.RFC3339, modTime)
	handlErr(err)

	err = os.Chtimes(fileName, aTime, mTime)
	handlErr(err)
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	return false
}

func ExecSystemCmd(operatingSystem string, arguments []string) string {
	var systemShell string
	var firstArg []string

	if operatingSystem == "unix" {
		systemShell = "bash"
		firstArg = append(firstArg, "-c")
	} else if operatingSystem == "windows" {
		systemShell = "cmd"
		firstArg = append(firstArg, "/C")
	} else {
		return "ERROR: NOT A KNOWN OPERATING SYSTEM"
	}

	firstArg = append(firstArg, arguments...)

	cmd := exec.Command(systemShell, firstArg...)

	combinedOut, err := cmd.CombinedOutput()
	handlErr(err)

	return string(combinedOut[:])
}
