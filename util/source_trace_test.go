package util

import "testing"

func TestGetLine(t *testing.T) {
	s:="/go/src/github.com/curtank/go-explicit-type-conversion/examples/user_info.go"
	r,_:=GetLine(s,95)
	t.Log(r)
	t.Error("todo")
}
func TestGetFileFromStack(t *testing.T){
	s:=`goroutine 1 [running]:
runtime/debug.Stack(0xc000053b20, 0x40bad8, 0x160)
        /usr/local/go/src/runtime/debug/stack.go:24 +0x9d
runtime/debug.PrintStack()
        /usr/local/go/src/runtime/debug/stack.go:16 +0x22
github.com/curtank/go-explicit-type-conversion/client.(*Client).AddFunc(0xc000053c60, 0x589940, 0x5c5118, 0xc000053c80, 0x4142fa)
        /go/src/github.com/curtank/go-explicit-type-conversion/client/conversion.go:190 +0x34
main.main()
        /go/src/github.com/curtank/go-explicit-type-conversion/examples/user_info.go:94 +0xd1

`	
	l,i,_:=GetFileFromStack(s)
	t.Log(l,i)
	t.Error("todo")
	}