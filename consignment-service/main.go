package main

import (
	"context" // proto编译出的包，使用到context
	"google.golang.org/grpc"
	"log"
	"net"
	pb "shippy/consignment-service/proto/consignment" // 为go_micro_srv_consignment包取个别名pb
)

const (
	PORT = ":50051"
)

// 仓库接口
// 用来处理一些业务逻辑。
type IRepository interface {
	// 存放新货物。
	// 后面需要把相关的数据存放到数据库。
	Create(consignment *pb.Consignment) (*pb.Consignment, error)

	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

// 实现接口
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	// 追加到consignments里面。
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// 定义微服务
type service struct {
	repo Repository
}

// 需要实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	resp := &pb.Response{Consignments: consignments}

	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen on :%s\n", PORT)

	// 创建 grpc server 实例。
	server := grpc.NewServer()

	// 向 gRPC 服务器注册微服务
	pb.RegisterShippingServiceServer(server, &service{})

	// 启动 gRPC 服务器
	// 阻塞等待用户调用
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
