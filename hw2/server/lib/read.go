package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
)

func ReadNFS(in *pb.ReadArgs) (*pb.ReadReturn, error) {

	// we're not using the syscalls directly here, so we're not dealing
	// directly with file descriptors
	// the reason is because we need to do fseek
	// if we wanted to cache open file descriptors, we would just have to do
	// that separately

	// get path for the file
	filepath, err := InumToPath(int(in.Fh.Inode))

	// TODO: what do we do with genum?
	// f_genum := in.Fh.Genum

	// open that file to get file object for it (this isn't an fd)
	f, err := os.Open(filepath)

	// use fseek to go to offset position in the file
	new_offset, err := f.Seek(in.Offset, 0)
	new_offset = new_offset // to supress compiler warning

	// use read to read count bytes
	f_count := in.Count
	data := make([]byte, f_count, f_count)
	n_bytes_read, err := f.Read(data)
	n_bytes_read = n_bytes_read // to supress compiler warning

	// TODO: if f_count is greater than the number of bytes to read in data, do
	// we need to do anything special?

	// TODO: what do we do with file attributes?

	return &pb.ReadReturn{Attr: &pb.Attribute{}, Data: data}, err
}
