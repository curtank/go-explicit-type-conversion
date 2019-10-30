package client

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"runtime/debug"
	"text/template"

	"github.com/curtank/go-explicit-type-conversion/util"
	// "github.com/ztrue/tracerr"
)

var (
	ErrNotFunc          = errors.New("not a func")
	ErrCallConvertFunc  = errors.New("call convert func with error return and but can not assert to error type")
	ErrNotStruct        = errors.New("convert type is not a struct")
	ErrFuncAlreadyExist = errors.New("adding func alreadyexist")
	ErrFuncNotExist     = errors.New("can not find func")
	ErrFuncInNotOne     = errors.New("func input len is not 1")
	ErrFuncOutNotTwo    = errors.New("func output len is not 2")
)
var (
	generate_tmpl = `
func f({{.InName}} {{.InType}}) ({{.OutType}},error){
	{{.OutName}}:={{.OutType}}{
		{{range .Direct}}{{.OutName}}:{{$.InName}}.{{.InName}},
		{{end}}
	}
	{{range .InDirect}}
	{{.OutName}},err:={{.FunctionName}}({{$.InName}}.{{.InName}})
	if err !=nil{
		return  {{$.OutName}},err
	}
	{{$.OutName}}.{{.OutName}}={{.OutName}}
	{{end}}
	return {{.OutName}},nil
}
`
	generate_tmp, _ = template.New("generate").Parse(generate_tmpl)
	removeMainReg   = regexp.MustCompile(`main[.](.+)`)
	funcnameReg=regexp.MustCompile(`AddFunc\((.+)\)`)
)

//Client Conversion Client
type Client struct {
	transMap transFuncMap
	staticCall staticCallMap 
}
type inOutPair struct {
	In  reflect.Type
	Out reflect.Type
}
type transFuncMap map[inOutPair]reflect.Value
type staticCallMap map[inOutPair]string

type directTrans struct {
	InName  string
	OutName string
}
type inDirectTrans struct {
	InName       string
	OutName      string
	FunctionName string
}
type transOp struct {
	Direct   []directTrans
	InDirect []inDirectTrans
	InName   string
	OutName  string
	OutType  string
	InType   string
}

func NewClient() *Client {
	t := transFuncMap{}
	s:=staticCallMap{}
	c := Client{transMap: t,staticCall: s}
	return &c
}

//Convert a sturct to another
func (c *Client) Convert(from interface{}, to interface{}) error {
	fromv := reflect.Indirect(reflect.ValueOf(from))
	tov := reflect.Indirect(reflect.ValueOf(to))
	if !IsStruct(fromv) || !IsStruct(tov) {
		return ErrNotStruct
	}
	for i := 0; i < fromv.NumField(); i++ {
		in := fromv.Field(i)
		correspondingField, ok := c.findCorrespondingField(fromv.Type().Field(i), &tov)
		// fmt.Fprintln(os.Stderr, correspondingField)
		if ok {
			out := tov.FieldByName(correspondingField.Name)
			if out.IsValid() {
				iop, err := c.ConvertField(in, out)
				// fmt.Fprintln(os.Stderr, iop, err)

				if err != nil {
					if err == ErrFuncNotExist {
						continue
					}
					return err
				}
				out.Set(iop)
			}
		}
	}
	return nil
}

//StaticGenerate a sturct to another
func (c *Client) StaticGenerate(from interface{}, to interface{}, inName, outName string) (string,error) {
	fromv := reflect.Indirect(reflect.ValueOf(from))
	tov := reflect.Indirect(reflect.ValueOf(to))
	if !IsStruct(fromv) || !IsStruct(tov) {
		return "",ErrNotStruct
	}
	transop := transOp{InName: inName, OutName: outName, OutType: tov.Type().Name(), InType: fromv.Type().Name()}

	for i := 0; i < fromv.NumField(); i++ {
		in := fromv.Field(i)
		correspondingField, ok := c.findCorrespondingField(fromv.Type().Field(i), &tov)
		// fmt.Fprintln(os.Stderr, correspondingField)
		// fmt.Fprintln(os.Stderr, fromv.Type().Field(i).Type)
		// fmt.Fprintln(os.Stderr, correspondingField.Type)

		if ok {
			out := tov.FieldByName(correspondingField.Name)
			if out.IsValid() {
				fname, code, err := c.StaticGenerateField(in, out)

				// fmt.Fprintln(os.Stderr, iop, err)

				if err != nil {
					if err == ErrFuncNotExist {
						continue
					}
					return "",err
				}
				if code == 0 {
					// fmt.Fprintln(os.Stderr, 0, fromv.Type().Field(i).Name, correspondingField.Name, fname)
					transop.Direct = append(transop.Direct, directTrans{InName: fromv.Type().Field(i).Name, OutName: correspondingField.Name})
				}
				if code == 1 {
					// fmt.Fprintln(os.Stderr, 1, fromv.Type().Field(i).Type, correspondingField.Type, fname)
					transop.InDirect = append(transop.InDirect, inDirectTrans{
						InName:       removeMain(fromv.Type().Field(i).Name),
						OutName:      removeMain(correspondingField.Name),
						FunctionName: fname,
					})
				}
				// out.Set(iop)
			}
		}
	}
	// fmt.Fprintln(os.Stderr, transop)

	buf := new(bytes.Buffer)

	generate_tmp.Execute(buf, transop)
	// fmt.Fprintln(os.Stderr, buf.String())

	return buf.String(),nil
}
func removeMain(s string) string {
	res := removeMainReg.ReplaceAllString(s, "${1}")

	return res
}
func (c *Client) findCorrespondingField(filed reflect.StructField, tostruct *reflect.Value) (reflect.StructField, bool) {
	return tostruct.Type().FieldByName(filed.Name)
}
func IsStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Struct
}
func (*Client) generatekey(o interface{}) {

}
func (c *Client) findFunc(in, out reflect.Type) reflect.Value {
	pp := inOutPair{In: in, Out: out}
	return c.transMap[pp]
}

func (c *Client) call(para reflect.Value, iop inOutPair) reflect.Value {
	nin := make([]reflect.Value, 1)
	nin[0] = reflect.ValueOf(para.Interface())
	return c.transMap[iop].Call(nin)[0]
}
func (c *Client) AddFunc(o interface{}) error {


	finterface := reflect.ValueOf(o)
	if finterface.Kind() != reflect.Func {
		return ErrNotFunc
	}
	functype := finterface.Type()
	if functype.NumIn() != 1 {
		return ErrFuncInNotOne
	}
	if functype.NumOut() != 2 {
		return ErrFuncOutNotTwo
	}
	in := functype.In(0)
	out := functype.Out(0)

	p := inOutPair{In: in, Out: out}
	if _, ok := c.transMap[p]; ok {
		return ErrFuncAlreadyExist
	}
	c.transMap[p] = finterface

	line,index,err:=util.GetFileFromStack(string(debug.Stack()))
	if err !=nil{
		return err
	}
	code,err:=util.GetLine(line,index)
	if err !=nil{
		return err
	}
	// fmt.Fprintln(os.Stderr, code)
	match:=funcnameReg.FindStringSubmatch(code)
	// fmt.Fprintln(os.Stderr, match[1])
	c.staticCall[p]=match[1]
	return nil
}
func (c *Client) ConvertField(in, out reflect.Value) (reflect.Value, error) {

	pair := inOutPair{In: in.Type(), Out: out.Type()}
	// fmt.Fprintln(os.Stderr, pair)

	f, ok := c.transMap[pair]
	if !ok {
		if in.Type() == out.Type() {
			return in, nil
		}
		return reflect.Value{}, ErrFuncNotExist
	}
	nin := make([]reflect.Value, 1)
	nin[0] = in
	res := f.Call(nin)
	// fmt.Fprintln(os.Stderr, res)

	err, e := res[1].Interface().(error)
	// fmt.Fprintln(os.Stderr, err, e)
	if e {
		return res[0], err
	}
	return res[0], err
}
func (c *Client) StaticGenerateField(in, out reflect.Value) (string, int8, error) {

	pair := inOutPair{In: in.Type(), Out: out.Type()}
	// fmt.Fprintln(os.Stderr, pair)

	f, ok := c.staticCall[pair]
	if !ok {
		if in.Type() == out.Type() {
			// fmt.Fprintln(os.Stderr, in.)
			return "", 0, nil
		}
		return "", -1, ErrFuncNotExist
	}
	return f, 1, nil
	// nin := make([]reflect.Value, 1)
	// nin[0] = in
	// res := f.Call(nin)
	// // fmt.Fprintln(os.Stderr, res)

	// err, e := res[1].Interface().(error)
	// // fmt.Fprintln(os.Stderr, err, e)
	// if e {
	// 	return res[0], err
	// }
	// return res[0], err
}
