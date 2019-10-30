package util

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
)

var (
	lineReg = regexp.MustCompile(`AddFunc[^\n]+\n[^\n]+\n[^\n]+\n[ \t]+(?P<File>[^ :]+):(?P<Line>[0-9]+)`)
)

func GetLine(location string, lineNum int64) (string, error) {
	f, err := os.Open(location)
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	sc := bufio.NewScanner(r)
	lastLine:= int64(0) 
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), sc.Err()
		}
	}
	return "", io.EOF
	defer f.Close()
	return "", nil
}
func GetFileFromStack(s string) (string, int64, error) {
	r := lineReg.FindStringSubmatch(s)
	// fmt.Fprintln(os.Stderr, reflect.TypeOf(r).String())
	// for i := range r {
	// 	fmt.Fprintln(os.Stderr, r)
	// }
	number,err:=strconv.ParseInt(r[2],10,16)
	
	if err !=nil{
		return "",-1,err
	}
	return r[1], number, nil
}
