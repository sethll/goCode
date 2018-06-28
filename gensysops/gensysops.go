/*
Package gensysops provides generic functions which you can expose to the otto
JavaScript runtime to perform common system operations. Examples include file
modification and even arbitrary system commands.

	import (
		"github.com/robertkrimen/otto"
		"github.com/sethll/goCode/gensysops"
	)

Use

To use a "gensysops" function, create an otto VM and Set the desired function
to a keyword.

	ottoVM := otto.New()
	ottoVM.Set("goFileExists", gensysops.FileExists)
	ottoVM.Run(`
		var myFile = "./test1.txt";
		var myFileExists = goFileExists(myFile);
		console.log(myFile, "exists:", myFileExists);
	`)

For explicit usage examples, refer to
https://github.com/sethll/goCode/gensysops/test.go
*/
package gensysops

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
