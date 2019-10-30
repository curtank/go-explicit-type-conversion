package examples

import "github.com/golang/protobuf/ptypes"

func f(user User) (UserInfo, error) {
	userinfo := UserInfo{
		Name:     user.Name,
		LastName: user.LastName,
	}

	Friends, err := UserIDs2Briefs(user.Friends)
	if err != nil {
		return userinfo, err
	}
	userinfo.Friends = Friends

	Follower, err := UserIDs2Briefs(user.Follower)
	if err != nil {
		return userinfo, err
	}
	userinfo.Follower = Follower

	RegisterTime, err := ptypes.TimestampProto(user.RegisterTime)
	if err != nil {
		return userinfo, err
	}
	userinfo.RegisterTime = RegisterTime

	LastLoginTime, err := github.com / golang / protobuf / ptypes.TimestampProto(user.LastLoginTime)
	if err != nil {
		return userinfo, err
	}
	userinfo.LastLoginTime = LastLoginTime

	return userinfo, nil
}
