package nfs

import (
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"syscall"
)

func StatToAttr(in *syscall.Stat_t) *pb.Attribute {
	var new_attr pb.Attribute
	new_attr.Uid = in.Uid
	new_attr.Gid = in.Gid
	new_attr.Mode = in.Mode
	new_attr.Size = uint64(in.Size)
	new_attr.Atime = uint64(in.Atim.Sec)
	new_attr.Mtime = uint64(in.Mtim.Sec)
	new_attr.AtimeNsec = uint32(in.Atim.Nsec)
	new_attr.MtimeNsec = uint32(in.Mtim.Nsec)
	return &new_attr
}
