package main

import (
	"fmt"
	"time"

	"github.com/curtank/go-explicit-type-conversion/client"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type UserInfo struct {
	Name         string
	Friends      []UserBrief
	RegisterTime *timestamp.Timestamp
}

type UserBrief struct {
	Name  string
	Phone string
}
type User struct {
	ID           string
	Name         string
	Friends      []string
	RegisterTime time.Time
}

func timetotimstamp(t time.Time) (*timestamp.Timestamp, error) {
	return ptypes.TimestampProto(t)
}
func UserIDs2Briefs(IDs []string) ([]UserBrief, error) {
	briefs := make([]UserBrief, len(IDs))
	for index, v := range IDs {
		briefs[index] = queryUser(v)
	}
	return briefs, nil
}
func queryUser(ID string) UserBrief {
	return UserBrief{
		Name: "Bob",
	}
}
func main() {
	c := client.NewClient()
	c.AddFunc(timetotimstamp)
	c.AddFunc(UserIDs2Briefs)
	user := User{
		ID:           "skIDsq",
		Name:         "Bill",
		Friends:      []string{"qqrrpp", "ssaaaee"},
		RegisterTime: time.Now(),
	}
	userinfo := UserInfo{}

	c.Convert(&user, &userinfo)
	fmt.Println(userinfo)
}
