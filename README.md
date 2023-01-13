# id_generator
### Usage
[1] update [zookeeper/setup.go](https://github.com/nova-plixx/id_generator/blob/main/zookeeper/setup.go) with desired
  - zookeeper ensemble connection string
  - shared int64 seed
  - start epoch milli
 
it defaults to 
  - zookeeper ensemble in localhost => *127.0.0.1:2181*
  - shared int64 seed => *9223372036854775783*
  - start epoch milli => *1577836800000* => *2020-01-01 00:00:00.000*
 
[2] run the zookeeper setup file => `go run zookeeper/setup.go`

[3] update [server/main.go](https://github.com/nova-plixx/id_generator/blob/main/server/main.go) with desired
  - zookeeper ensemble connection string

[4] run the server file => `go run server/main.go`

[5] sample client code can be found at [example-client/main.go](https://github.com/nova-plixx/id_generator/blob/main/example-client/main.go)
