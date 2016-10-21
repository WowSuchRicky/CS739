package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

// we're not using the syscalls directly here, so we're not dealing
// directly with file descriptors
// the reason is because we need to do fseek
// if we wanted to cache open file descriptors, we would just have to do
// that separately

func ReadNFS(in *pb.ReadArgs) (*pb.ReadReturn, error) {
	filepath, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.ReadReturn{}, errors.New("inode not found")
	}

	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(filepath)
	if err != nil {
		fmt.Println("Genum of file not found, fatal error.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.ReadReturn{}, errors.New("genum mismatch")
	}

	// open that file to get file object for it (this isn't an fd)
	f, err := os.Open(filepath)
	if err != nil { // propagate the file open error back to client
		fmt.Println("os.Open returned err")
		return &pb.ReadReturn{}, err
	}

	// use fseek to go to offset position in the file
	new_offset, err := f.Seek(in.Offset, 0)
	if err != nil { // propagate fseek error back to client
		fmt.Println("Fseek returned error")
		return &pb.ReadReturn{}, err
	}
	new_offset = new_offset // to supress compiler warning

	// use read to read count bytes
	f_count := in.Count
	data := make([]byte, f_count, f_count)
	n_bytes_read, err := f.Read(data)
	if err != nil {
		fmt.Println("Read returned error")
		return &pb.ReadReturn{}, err
	}

	n_bytes_read = n_bytes_read // to supress compiler warning

	// TODO: if f_count is greater than the number of bytes to read in data, do
	// we need to do anything special?

	// stat file to fill attributes
	var f_info syscall.Stat_t
	err = syscall.Stat(filepath, &f_info)
	if err != nil {
		fmt.Println("Stat failed, FATAL error")
		os.Exit(-1)
	}

	return &pb.ReadReturn{Attr: StatToAttr(&f_info), Data: data}, nil
}
