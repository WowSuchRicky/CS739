package nfs

import (
	pb "../../protos"
)

func ReadNFS(in *pb.ReadArgs) (*pb.ReadReturn, error) {
	return &pb.ReadReturn{Attr: &pb.Attribute{}, Data: []byte{1, 2}}, nil
}
