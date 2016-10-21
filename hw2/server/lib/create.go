package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

// two approaches:
// 1) get fd for the dir, then use openat() to create file in it
// 2) use open() or creat() straight-up with full path
// we'll take the 2nd approach for now

func CreateNFS(in *pb.CreateArgs) (*pb.CreateReturn, error) {
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

	// assume that in.Attr.Mode exists
	new_fd, err := syscall.Creat(new_file_path, in.Attr.Mode)
	if err != nil {
		fmt.Printf("Couldn't create file, error (check mode correct and file doesn't exist)")
		return &pb.CreateReturn{}, errors.New("could not create")
	}

	err = syscall.Close(new_fd)

	// get inode and genum of new file
	var new_f_info syscall.Stat_t
	err = syscall.Stat(new_file_path, &new_f_info)
	new_inode := new_f_info.Ino
	new_genum, err := PathToGen(new_file_path)

	return &pb.CreateReturn{Newfh: &pb.FileHandle{Inode: new_inode, Genum: new_genum}, Attr: StatToAttr(&new_f_info)}, err
}
