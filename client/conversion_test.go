package client

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func TestPointerBehaivor4(t *testing.T) {
	type NameStr struct {
		Name string
	}
	type OA struct {
		Name *NameStr
		ID   string
	}
	type OB struct {
		Name string
		ID   string
	}
	f := func(n *NameStr) (string, error) {
		return n.Name, nil
	}
	// f2 := func(n *NameStr) (string, error) {
	// 	return n.Name, errors.New("convert  failed")
	// }
	n := NameStr{Name: "sss"}
	oa := OA{Name: &n, ID: "2212"}
	ob := OB{}
	c := NewClient()
	err := c.AddFunc(f)
	t.Log(err)
	t.Log(c.transMap)
	err = c.Convert(&oa, &ob)
	t.Log(err)
	t.Log(ob)
	t.Error("todo")
}
func TestNotStruct(t *testing.T) {
	c := NewClient()
	err := c.Convert(12, 22)
	if err != ErrNotStruct {
		t.Error("not struct failed")
	}
}
func SuccessConvert(t *testing.T, from, to interface{}, c *Client) {
	err := c.Convert(from, to)
	if err != nil {
		t.Error(err)
	}
}
func SuccessAdd(t *testing.T, f interface{}, c *Client) {
	err := c.AddFunc(f)
	if err != nil {
		t.Error(err)
	}
}
func TestConvert1(t *testing.T) {
	type S1 struct {
		Name string
	}
	type S2 struct {
		Name string
	}
	c := NewClient()
	s1 := S1{Name: "hello"}
	s2 := S2{}
	SuccessConvert(t, &s1, &s2, c)
	if s2.Name != "hello" {
		t.Error("fail convert", s2)
	}
}

func TestConvert2(t *testing.T) {
	type SS struct {
		Value string
	}
	type S1 struct {
		Name string
	}
	type S2 struct {
		Name *SS
	}
	c := NewClient()
	s1 := S1{Name: "hello"}
	s2 := S2{}
	s2SS := func(s string) (*SS, error) {
		return &SS{Value: s}, nil
	}
	SuccessAdd(t, s2SS, c)
	SuccessConvert(t, &s1, &s2, c)
	if s2.Name.Value != "hello" {
		t.Error("fail convert", s2)
	}
}
func TestConvert3(t *testing.T) {
	type S1 struct {
		Name string
	}
	type S2 struct {
		Name string
	}
	c := NewClient()
	s1 := S1{Name: "hello"}
	s2 := S2{}
	s2SS := func(s string) (string, error) {
		return s + " world", nil
	}
	SuccessAdd(t, s2SS, c)
	SuccessConvert(t, &s1, &s2, c)
	if s2.Name != "hello world" {
		t.Error("fail convert", s2)
	}
}
func TestTime(t *testing.T) {

	type GoTimeStamp struct {
		CreateTime time.Time
		EndTime    time.Time
	}
	type GRPCTimeStamp struct {
		CreateTime *timestamp.Timestamp
		EndTime    *timestamp.Timestamp
	}

	f := func(t time.Time) (*timestamp.Timestamp, error) {
		return ptypes.TimestampProto(t)

	}
	c := NewClient()
	SuccessAdd(t, f, c)
	gotime := GoTimeStamp{CreateTime: time.Now(), EndTime: time.Now()}
	grpctime := GRPCTimeStamp{}
	SuccessConvert(t, &gotime, &grpctime, c)
	if grpctime.CreateTime.Seconds != gotime.CreateTime.Unix() {
		t.Error("fail convert", gotime, grpctime)
	}
}

// func TestConvertFiled(t *testing.T) {
// 	type Bob struct {
// 		Name string
// 	}
// 	type Alice struct {
// 		Name int
// 	}
// 	f := s2i
// 	c := NewClient()
// 	res := c.AddFunc(f)
// 	// c.ConvertField()
// }
