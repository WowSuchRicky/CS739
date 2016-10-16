package nfs

import (
	pb "../../protos"
)

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {
	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: 1, Fsnum: 2, Genum: 32}, Attr: &pb.Attribute{}}, nil
}
