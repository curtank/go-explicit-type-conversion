# go-explicit-type-conversion
explicit type conversion from one type to another

suppose you have 2 type of time the `time.Time` and `timestamp.Timestamp`

```go
import (
	"fmt"
	"time"

	"github.com/curtank/go-explicit-type-conversion/client"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type GoTimeStamp struct {
	CreateTime time.Time
	EndTime    time.Time
}
type GRPCTimeStamp struct {
	CreateTime *timestamp.Timestamp
	EndTime    *timestamp.Timestamp
}

func timetotimstamp(t time.Time) (*timestamp.Timestamp, error) {
	return ptypes.TimestampProto(t)
}
```

```go
func main() {
	c := client.NewClient()
	c.AddFunc(timetotimstamp)
	gotime := GoTimeStamp{CreateTime: time.Now(), EndTime: time.Now()}
	grpctime := GRPCTimeStamp{}
	c.Convert(&gotime, &grpctime)
	fmt.Println(grpctime)
}


```
