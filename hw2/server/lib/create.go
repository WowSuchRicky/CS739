package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {

	// get path of directory
	dir_path, _ := InumToPath(int(in.Dirfh.Inode))

	// get full path
	full_path := dir_path + "/" + in.Name

	// two approaches:
	// 1) get fd for the dir, then use openat() to create file in it
	// 2) use open() or creat() straight-up with full path
	// we'll take the 2nd approach for now

	// TODO: how should we use attributes provided? (maybe don't need to)
	// I think attributes in this case should just be the mode (i.e. 00666 as we use below?)
	// file_attr := in.Attr

	new_fd, err := syscall.Creat(full_path, 00666)
	err = syscall.Close(new_fd)
	// TODO: should we close this immediately as we're doing above?

	// get inode and genum of new file
	var new_f_info syscall.Stat_t
	err = syscall.Stat(full_path, &new_f_info)
	new_inode := new_f_info.Ino
	new_genum := uint64(1) // TODO: get genum

	// what attributes do we return; probably just use the new_f_info from above to populate it

	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: new_inode, Genum: new_genum}, Attr: &pb.Attribute{}}, err
}
