package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {

	// get path of directory
	dir_path, _ := InumToPath(int(in.Dirfh.Inode))

	// TODO: how should do we use input genum?
	// dir_genum := in.Dirfh.Genum

	// get full path
	full_path := dir_path + "/" + in.Name

	// two approaches:
	// 1) get fd for the dir, then use openat() to create file in it
	// 2) use open() or creat() straight-up with full path
	// we'll take the 2nd approach for now

	// file_attr := in.Attr // TODO: how should we use attributes provided?
	new_fd, err := syscall.Creat(full_path, 0)
	err = syscall.Close(new_fd) // TODO: should we close this immediately?

	// get inode and genum of new file
	var new_f_info syscall.Stat_t
	err = syscall.Stat(full_path, &new_f_info)
	new_inode := new_f_info.Ino
	new_genum := uint64(1) // TODO: get genum

	// TODO: what attributes do we return?

	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: new_inode, Genum: new_genum}, Attr: &pb.Attribute{}}, err
}
