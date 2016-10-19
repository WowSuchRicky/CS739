package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func LookupNFS(in *pb.LookupArgs) (*pb.LookupReturn, error) {

	// 1) get path of directory using inode & concatenate it with file name
	dir_path, _ := InumToPath(int(in.Dirfh.Inode))
	full_path := dir_path + in.Name // TODO: do we need to add forward slash?

	// 2) get inode & genum of that file
	var f_info syscall.Stat_t
	syscall.Stat(full_path, &f_info)
	ret_inode := int32(f_info.Ino)
	ret_genum := int32(1) // TODO: genum

	// 5) TODO: get attributes

	return &pb.LookupReturn{Fh: &pb.FileHandle{Inode: ret_inode, Genum: ret_genum}, Attr: &pb.Attribute{}}, nil
}
