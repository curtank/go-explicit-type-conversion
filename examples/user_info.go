package main

import (
	"fmt"
	"time"

	"github.com/curtank/go-explicit-type-conversion/client"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type UserInfo struct {
	Name          string
	Friends       []UserBrief
	Follower      []UserBrief
	RegisterTime  *timestamp.Timestamp
	LastLoginTime *timestamp.Timestamp
}

type UserBrief struct {
	Name  string
	Phone string
}
type User struct {
	ID            string
	Name          string
	Friends       []string
	Follower      []string
	RegisterTime  time.Time
	LastLoginTime time.Time
}

var storage = map[string]string{
	"1": "Sir Humphrey Appleby",
	"2": "Bernard Woolley",
	"3": "Jim Hacker",
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
		Name: storage[ID],
	}
}
func main() {
	c := client.NewClient()
	c.AddFunc(timetotimstamp)
	c.AddFunc(UserIDs2Briefs)
	user := User{
		ID:            "1",
		Name:          "Sir Humphrey Appleby",
		Friends:       []string{"2", "3"},
		Follower:      []string{"2"},
		RegisterTime:  time.Now(),
		LastLoginTime: time.Now(),
	}
	userinfo := UserInfo{}

	c.Convert(&user, &userinfo)
	fmt.Println(userinfo)
}
