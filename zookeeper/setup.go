package main

import (
	"github.com/go-zookeeper/zk"
	"math"
	"strconv"
	"time"
)

func main() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, _ = conn.Create("/id_generator", nil, 0, zk.WorldACL(zk.PermAll))
	_, _ = conn.Create("/id_generator/shared_seed", []byte(strconv.Itoa(9223372036854775783)), 0, zk.WorldACL(zk.PermAll))
	_, _ = conn.Create("/id_generator/start_epoch", []byte(strconv.Itoa(1577836800000)), 0, zk.WorldACL(zk.PermAll))
	_, _ = conn.Create("/id_generator/instance_id", nil, 0, zk.WorldACL(zk.PermAll))

	const uniquePartBits = 10
	var maxInstanceNumber = int(math.Ceil(math.Pow(2, float64(uniquePartBits))))
	for i := 0; i < maxInstanceNumber; i++ {
		instanceId := "/id_generator/instance_id/" + strconv.Itoa(i)
		_, err = conn.Create(instanceId, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			panic(err)
		}
	}
}
