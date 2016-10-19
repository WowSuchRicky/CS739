package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
)

func WriteNFS(in *pb.WriteArgs) (*pb.WriteReturn, error) {
	// we have same input args as Read but also have Data which is an array of bytes
	// see read.go for info on some decision made here

	// get path for the file
	filepath, err := InumToPath(int(in.Fh.Inode))

	// TODO: what do we do with genum?
	// f_genum := in.Fh.Genum

	// open that file to get file object for it (this isn't an fd)
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0)

	// write the data into it starting at in.Offset
	data_to_write := in.Data[0:in.Count]
	n_bytes_written, err := f.WriteAt(data_to_write, in.Offset)
	n_bytes_written = n_bytes_written // to supress compiler warning

	// TODO: pull out attributes of that file AFTER writing it

	return &pb.WriteReturn{Attr: &pb.Attribute{}}, err
}
