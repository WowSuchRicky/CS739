package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {

	// get path of directory
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	new_file_path := dir_path + "/" + in.Name

	// error if inode wasn't found
	if err != nil {
		return &pb.CreateReturn{}, errors.New("inode not found")
	}

	// error if genum not found (fatal error, should never happen), or genum mismatch
	fh_genum := in.Dirfh.Genum
	fs_genum, err := PathToGen(dir_path)
	if err != nil {
		fmt.Println("Genum not found, error.")
		os.Exit(-1)
	}
	if fh_genum != fs_genum {
		fmt.Printf("fh_genum: %v, fs_genum: %v", fh_genum, fs_genum)
		return &pb.CreateReturn{}, errors.New("genum mismatch")
	}

	// two approaches:
	// 1) get fd for the dir, then use openat() to create file in it
	// 2) use open() or creat() straight-up with full path
	// we'll take the 2nd approach for now

	// TODO: how should we use attributes provided? (maybe don't need to)
	// I think attributes in this case should just be the mode (i.e. 00666 as we use below?)
	// file_attr := in.Attr

	new_fd, err := syscall.Creat(new_file_path, 00666)
	err = syscall.Close(new_fd)

	// get inode and genum of new file
	var new_f_info syscall.Stat_t
	err = syscall.Stat(new_file_path, &new_f_info)
	new_inode := new_f_info.Ino
	new_genum, err := PathToGen(new_file_path)

	// what attributes do we return; probably just use the new_f_info from above to populate it

	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: new_inode, Genum: new_genum}, Attr: &pb.Attribute{}}, err
}
