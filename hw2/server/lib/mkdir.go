package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func MkdirNFS(in *pb.MkdirArgs) (*pb.MkdirReturn, error) {
	dir_path, err := InumToPath(int(in.Dirfh.Inode))
	new_dir_path := dir_path + "/" + in.Name

	// error if inode wasn't found
	if err != nil {
		return &pb.MkdirReturn{}, errors.New("inode not found")
	}

	// error if genum not found (fatal error, should never happen), or genum mismatch
	fh_genum := in.Dirfh.Genum
	fs_genum, err := PathToGen(dir_path)
	if err != nil {
		fmt.Println("Genum not found, fatal error.")
		os.Exit(-1)
	}
	if fh_genum != fs_genum {
		fmt.Printf("fh_genum: %v, fs_genum: %v", fh_genum, fs_genum)
		return &pb.MkdirReturn{}, errors.New("genum mismatch")
	}

	err = syscall.Mkdir(new_dir_path, in.Attr.Mode)
	if err != nil {
		fmt.Printf("Couldn't create directory")
		return &pb.MkdirReturn{}, errors.New("could not create")
	}

	// ge filehandle for the new directory
	var new_dir_info syscall.Stat_t
	err = syscall.Stat(new_dir_path, &new_dir_info)
	new_inode := new_dir_info.Ino
	new_genum, err := PathToGen(new_dir_path)

	if err != nil {
		fmt.Println("Genum not found, fatal error.")
		os.Exit(-1)
	}

	return &pb.MkdirReturn{Fh: &pb.FileHandle{Inode: new_inode, Genum: new_genum}, Attr: StatToAttr(&new_dir_info)}, nil
}
