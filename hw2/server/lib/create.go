package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {
	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: 1, Genum: 32}, Attr: &pb.Attribute{}}, nil
}
