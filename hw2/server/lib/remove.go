package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func RemoveNFS(in *pb.RemoveArgs) (*pb.RemoveReturn, error) {
	return &pb.RemoveReturn{Status: 52}, nil
}
