package nfs

import (
	pb "../../protos"
)

func WriteNFS(in *pb.WriteArgs) (*pb.WriteReturn, error) {
	return &pb.WriteReturn{Attr: &pb.Attribute{}}, nil
}
