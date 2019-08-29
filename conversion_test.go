package conversion

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

type A struct {
	Name *NameS `json:"name"`
}
type B struct {
	Name *NameS
}
type NameS struct {
	N string
}
type NameI struct {
	N int
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
func TestConversion(t *testing.T) {

	ns := NameS{N: "qqq"}
	a := A{Name: &ns}
	b := B(a)

	t.Error(b.Name.N)
	ns.N = "ppp"
	t.Error(b.Name.N)
}

type A2 struct {
	Name NameS
}
type B2 struct {
	Name NameS
}

func TestConversion2(t *testing.T) {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Trace.Println("kjsdflkjskdl")
	ns := NameS{N: "qqq"}
	a := A{Name: &ns}
	b := B(a)
	t.Log("sjkdflk")
	t.Log(b)
	// t.Error(b.Name.N)
	ns.N = "ppp"
	t.Error("TODO")
}

func TestIsStruct(t *testing.T) {
	ns := NameS{N: "qqq"}
	// res := IsStruct(&ns)
	// res := reflect.Indirect(reflect.ValueOf(&ns)).Kind()
	// t.Error(res)
	ps := reflect.ValueOf(&ns)
	// struct
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// exported field
		f := s.FieldByName("N")
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				// change value of N
				if f.Kind() == reflect.String {
					x := "rrr"

					f.SetString(x)
				}
			}
		}
	}
	// N at end
	t.Error(ns.N)
}
func i2i(a int) int {
	return a + 1
}
func s2i(s string) int {
	return len(s)
}
func TestFunc(t *testing.T) {
	funcs := make([]interface{}, 3, 3)               // I use interface{} to allow any kind of func
	funcs[0] = i2i                                   // good
	funcs[1] = func(a string) int { return len(a) }  // good
	funcs[2] = func(a string) string { return ":(" } // bad
	cur := NameS{N: "ss"}
	curV := reflect.ValueOf(cur)

	for _, fi := range funcs {
		f := reflect.ValueOf(fi)
		functype := f.Type()
		good := false
		for i := 0; i < functype.NumIn(); i++ {
			if "int" == functype.In(i).String() {
				good = true // yes, there is an int among inputs
				break
			}
			if "string" == functype.In(i).String() {
				in := make([]reflect.Value, 1)
				in[0] = reflect.ValueOf(curV.FieldByName("N").Interface())
				t.Error(f.Call(in)[0])
				good = true // yes, there is an int among inputs
				break
			}
		}
		for i := 0; i < functype.NumOut(); i++ {
			if "int" == functype.Out(i).String() {

				good = true // yes, there is an int among outputs
				break
			}
		}
		if good {

		}
	}

}
func TestAddFunc2(t *testing.T) {
	f := s2i
	finterface := reflect.ValueOf(f)
	if finterface.Kind() != reflect.Func {
		t.Error("jskldfj")
	}
	functype := finterface.Type()
	in := functype.In(0)
	out := functype.Out(0)

	p := inOutPair{In: in, Out: out}
	trans := transFuncMap{}
	trans[p] = finterface

	pp := inOutPair{In: in, Out: out}

	cur := NameS{N: "ss"}
	curV := reflect.ValueOf(cur)
	nin := make([]reflect.Value, 1)
	nin[0] = reflect.ValueOf(curV.FieldByName("N").Interface())

	t.Log(trans[pp].Call(nin)[0])
	t.Error("TODO")
}
func TestAddFunc(t *testing.T) {
	f := s2i
	f2 := func() string { return ":(" }
	c := NewClient()
	res := c.Addfunc(f)
	res = c.Addfunc(f)
	t.Log(res)
	res = c.Addfunc(f2)
	t.Log(res)
	functype := reflect.ValueOf(f).Type()
	in := functype.In(0)
	out := functype.Out(0)
	pp := inOutPair{In: in, Out: out}
	cur := NameS{N: "ss"}
	curV := reflect.ValueOf(cur)
	ret := c.call(curV.FieldByName("N"), pp)
	t.Log(ret)
	t.Error("TODO")
}
func prepareClient(t *testing.T) *Client {
	c := NewClient()
	c.Addfunc(s2i)
	return c
}
func TestConvert(t *testing.T) {
	ns := NameS{N: "qqq"}
	ni := NameI{}
	c := prepareClient(t)
	// fromv := reflect.Indirect(reflect.ValueOf(&ns))
	// tov := reflect.Indirect(reflect.ValueOf(&ni))
	// for i := 0; i < fromv.NumField(); i++ {
	// 	in := fromv.Field(i)
	// 	correspondingField, ok := c.findCorrespondingField(fromv.Type().Field(i), &tov)
	// 	if ok {
	// 		out := tov.FieldByName(correspondingField.Name)
	// 		if out.IsValid() {
	// 			iop := inOutPair{In: in.Type(), Out: out.Type()}
	// 			out.Set(c.call(in, iop))
	// 		}
	// 	}

	// }
	c.Convert(&ns, &ni)
	t.Log(ni)
	t.Error("TODO")
}

type GoTimeStamp struct {
	T time.Time
}
type GRPCTimeStamp struct {
	T *timestamp.Timestamp
}

func timetotimstamp(t time.Time) (*timestamp.Timestamp, error) {
	c, err := ptypes.TimestampProto(t)
	return c, err
}
func TestConvertTime(t *testing.T) {
	from := GoTimeStamp{T: time.Now()}
	to := GRPCTimeStamp{}
	Convey("Convert time", t, func() {
		c := NewClient()
		err := c.Addfunc(timetotimstamp)
		So(err, ShouldBeNil)
		err = c.Convert(&from, &to)
		So(err, ShouldBeNil)
		t.Log(to)
	})

}
func TestConvertPaste(t *testing.T) {
	type OA struct {
		Name string
		ID   string
	}
	type OB struct {
		Name int
		ID   string
	}
	oa := OA{Name: "bob", ID: "2212"}
	ob := OB{}
	c := NewClient()
	c.Addfunc(s2i)
	c.Convert(&oa, &ob)
	t.Log(ob)

	t.Error("todo")
}
func TestPointerBehaivor(t *testing.T) {
	type NameStr struct {
		Name string
	}
	type OA struct {
		Name NameStr
		ID   string
	}
	type OB struct {
		Name NameStr
		ID   string
	}
	n := NameStr{Name: "sss"}
	oa := OA{Name: n, ID: "2212"}
	ob := OB{}
	c := NewClient()
	c.Addfunc(s2i)
	c.Convert(&oa, &ob)
	t.Log(ob)
	t.Error("todo")
}
func TestPointerBehaivor2(t *testing.T) {
	type NameStr struct {
		Name string
	}
	type OA struct {
		Name *NameStr
		ID   string
	}
	type OB struct {
		Name *NameStr
		ID   string
	}
	Convey("pointer paste", t, func() {
		n := NameStr{Name: "sss"}
		oa := OA{Name: &n, ID: "2212"}
		ob := OB{}
		c := NewClient()
		c.Addfunc(s2i)
		r := c.Convert(&oa, &ob)
		So(r, ShouldBeNil)
		So(ob.Name.Name, ShouldEqual, "sss")
		n.Name = "qqq"
		So(ob.Name.Name, ShouldEqual, "qqq")
	})

}
func TestPointerBehaivor3(t *testing.T) {
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
	f := func(n *NameStr) string {
		return n.Name
	}
	n := NameStr{Name: "sss"}
	oa := OA{Name: &n, ID: "2212"}
	ob := OB{}
	c := NewClient()
	c.Addfunc(f)
	c.Convert(&oa, &ob)
	t.Log(ob)
	t.Error("todo")
}
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
	err := c.Addfunc(f)
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

// func TestConvertFiled(t *testing.T) {
// 	type Bob struct {
// 		Name string
// 	}
// 	type Alice struct {
// 		Name int
// 	}
// 	f := s2i
// 	c := NewClient()
// 	res := c.Addfunc(f)
// 	// c.ConvertField()
// }
