package conversion

import (
	"errors"
	"reflect"
)

var (
	ErrNotFunc          = errors.New("not a func")
	ErrNotStruct        = errors.New("convert type is not a struct")
	ErrFuncAlreadyExist = errors.New("adding func alreadyexist")
	ErrFuncNotExist     = errors.New("can not find func")
	ErrFuncInNotOne     = errors.New("func input len is not 1")
	ErrFuncOutNotOne    = errors.New("func output len is not 1")
)

//Client Conversion Client
type Client struct {
	transMap transFuncMap
}
type inOutPair struct {
	In  reflect.Type
	Out reflect.Type
}
type transFuncMap map[inOutPair]reflect.Value

func NewClient() *Client {
	t := transFuncMap{}
	c := Client{transMap: t}
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
		if ok {
			out := tov.FieldByName(correspondingField.Name)
			if out.IsValid() {
				iop, err := c.ConvertField(in, out)
				if err != nil {
					return err
				}
				out.Set(iop)
			}
		}
	}
	return nil
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
func (c *Client) Addfunc(o interface{}) error {
	finterface := reflect.ValueOf(o)
	if finterface.Kind() != reflect.Func {
		return ErrNotFunc
	}
	functype := finterface.Type()
	if functype.NumIn() != 1 {
		return ErrFuncInNotOne
	}
	if functype.NumOut() != 1 {
		return ErrFuncOutNotOne
	}
	in := functype.In(0)
	out := functype.Out(0)

	p := inOutPair{In: in, Out: out}
	if _, ok := c.transMap[p]; ok {
		return ErrFuncAlreadyExist
	}
	c.transMap[p] = finterface
	return nil
}
func (c *Client) ConvertField(in, out reflect.Value) (reflect.Value, error) {
	pair := inOutPair{In: in.Type(), Out: out.Type()}
	f, ok := c.transMap[pair]
	if !ok {
		return reflect.Value{}, ErrNotFunc
	}
	nin := make([]reflect.Value, 1)
	nin[0] = in
	res := f.Call(nin)[0]
	return res, nil
}
