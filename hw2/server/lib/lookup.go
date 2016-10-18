package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
)

func LookupNFS(in *pb.LookupArgs) (*pb.LookupReturn, error) {
	// approach #1: directly use the fs of the server
	// 1) change directory to input filehandle
	// 2) use os.Stat to get the inode number for the named file
	// 3) revert back to original directory

	// problems
	// - how would we handle gen nums? need to store separately? underlying file system doesn't necessarily this
	return &pb.LookupReturn{Fh: &pb.FileHandle{Inode: 1, Fsnum: 2, Genum: 32}, Attr: &pb.Attribute{}}, nil
}
