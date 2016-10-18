package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func ReadNFS(in *pb.ReadArgs) (*pb.ReadReturn, error) {
	return &pb.ReadReturn{Attr: &pb.Attribute{}, Data: []byte{1, 2}}, nil
}
