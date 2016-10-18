package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func WriteNFS(in *pb.WriteArgs) (*pb.WriteReturn, error) {
	return &pb.WriteReturn{Attr: &pb.Attribute{}}, nil
}
