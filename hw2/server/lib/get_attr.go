package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func GetAttrNFS(in *pb.GetAttrArgs, wq *ServerWriteQueue) (*pb.GetAttrReturn, error) {

	// get path for file
	file_path, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.GetAttrReturn{}, errors.New("inode not found")
	}

	// confirm genums are same
	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(file_path)
	if err != nil {
		fmt.Println("Genum not found, fatal err.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.GetAttrReturn{}, errors.New("genum mismatch")
	}

	// stat the file
	var f_info syscall.Stat_t
	err = syscall.Stat(file_path, &f_info)
	if err != nil {
		return &pb.GetAttrReturn{}, errors.New("stat failed")
	}

	// TODO: read through write buffer; udpate f_info file size with
	// max(current, offset + count from write)
	max_size := f_info.Size
	for i := 0; i < wq.size; i++ {
		write_inode := wq.queue[i].args.Fh.Inode
		if write_inode != in.Fh.Inode {
			continue
		}
		new_size := wq.queue[i].args.Offset + wq.queue[i].args.Count
		if new_size > max_size {
			max_size = new_size
		}
	}

	f_info.Size = max_size

	return &pb.GetAttrReturn{Attr: StatToAttr(&f_info)}, err
}
