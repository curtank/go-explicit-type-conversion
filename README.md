# go-explicit-type-conversion
explicit type conversion from one type to another

suppose you have 2 type of time the `time.Time` and `timestamp.Timestamp`

```go
type GoTimeStamp struct {
    CreateTime time.Time
    EndTime time.Time
}
type GRPCTimeStamp struct {
    CreateTime *timestamp.Timestamp
    EndTime *timestamp.Timestamp
}

func timetotimstamp(t time.Time) *timestamp.Timestamp {
    c, _ := ptypes.TimestampProto(t)
    return c
}

```
