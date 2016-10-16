package nfs

import (
	pb "../../protos"
)

func RemoveNFS(in *pb.RemoveArgs) (*pb.RemoveReturn, error) {
	return &pb.RemoveReturn{Status: 52}, nil
}
