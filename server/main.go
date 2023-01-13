package main

import (
	"context"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"plixx.dev/id_generator/domain"
	pb "plixx.dev/id_generator/proto"
	"strconv"
	"time"
)

type server struct {
	pb.UnimplementedIdServiceServer
	idGenerator domain.IdGenerator
}

func interceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if ctx.Err() != nil {
		freeZkGeneratorId()
		return nil, status.Error(codes.Canceled, "server is closing")
	}
	return handler(ctx, req)
}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	panicIfErrored(err)
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	pb.RegisterIdServiceServer(srv, &server{idGenerator: provideIdGenerator()})
	reflection.Register(srv)
	if e := srv.Serve(listener); e != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GenerateId(
	_ context.Context,
	_ *pb.GenerateIdRequest,
) (*pb.GenerateIdResponse, error) {
	ids, err := s.idGenerator.Generate(1)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateIdResponse{Id: ids[0]}, nil
}

func (s *server) GenerateMultipleIds(
	_ context.Context,
	request *pb.GenerateMultipleIdsRequest,
) (*pb.GenerateMultipleIdsResponse, error) {
	ids, err := s.idGenerator.Generate(int(request.Count))
	if err != nil {
		return &pb.GenerateMultipleIdsResponse{Ids: []int64{}}, err
	}
	return &pb.GenerateMultipleIdsResponse{Ids: ids}, nil
}

var generatorId int

func provideIdGenerator() domain.IdGenerator {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	panicIfErrored(err)
	defer conn.Close()

	sharedSeed := fetchZkInt64ConfigVal(conn, "/id_generator/shared_seed")
	startEpoch := fetchZkInt64ConfigVal(conn, "/id_generator/start_epoch")
	generatorId = fetchZkGeneratorId(conn)

	gen, err := domain.NewFinchIdGenerator(
		sharedSeed,
		generatorId,
		startEpoch,
		&domain.SystemClock{},
	)
	panicIfErrored(err)
	return gen
}

func fetchZkInt64ConfigVal(conn *zk.Conn, path string) int64 {
	valBytes, _, err := conn.Get(path)
	panicIfErrored(err)
	val, err := strconv.Atoi(string(valBytes))
	panicIfErrored(err)
	return int64(val)
}

func fetchZkGeneratorId(conn *zk.Conn) int {
	lock := acquireLock(conn, "/id_generator/instance_id/assign_lock")
	defer releaseLock(lock)
	children, _, err := conn.Children("/id_generator/instance_id")
	panicIfErrored(err)
	assignedId := children[0]
	err = conn.Delete("/id_generator/instance_id/"+assignedId, -1)
	panicIfErrored(err)
	id, err := strconv.Atoi(assignedId)
	panicIfErrored(err)
	log.Printf("Using generator id: %v", id)
	return id
}

func freeZkGeneratorId() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	panicIfErrored(err)
	defer conn.Close()
	lock := acquireLock(conn, "/id_generator/instance_id/release_lock")
	defer releaseLock(lock)
	instanceId := "/id_generator/instance_id/" + strconv.Itoa(generatorId)
	_, err = conn.Create(instanceId, nil, 0, zk.WorldACL(zk.PermAll))
	panicIfErrored(err)
	log.Printf("Releasing generator id: %v", generatorId)
}

func acquireLock(conn *zk.Conn, path string) *zk.Lock {
	lock := zk.NewLock(conn, path, zk.WorldACL(zk.PermAll))
	err := lock.Lock()
	panicIfErrored(err)
	return lock
}

func releaseLock(lock *zk.Lock) {
	err := lock.Unlock()
	panicIfErrored(err)
}

func panicIfErrored(err error) {
	if err != nil {
		panic(err)
	}
}
