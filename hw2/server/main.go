/*
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	nfs "github.com/Ricky54326/CS739/hw2/server/lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	port                 = ":50051"
	writeverf3_file_path = ".writeverf3"
)

var wq *(nfs.ServerWriteQueue)

// we implement NFSserver using server
type server struct{}

func (s *server) Lookup(ctx context.Context, in *pb.LookupArgs) (*pb.LookupReturn, error) {
	return nfs.LookupNFS(in)
}

func (s *server) Create(ctx context.Context, in *pb.CreateArgs) (*pb.CreateReturn, error) {
	return nfs.CreateNFS(in)
}

func (s *server) Remove(ctx context.Context, in *pb.RemoveArgs) (*pb.RemoveReturn, error) {
	return nfs.RemoveNFS(in)
}

func (s *server) Read(ctx context.Context, in *pb.ReadArgs) (*pb.ReadReturn, error) {
	return nfs.ReadNFS(in, wq)
}

func (s *server) Write(ctx context.Context, in *pb.WriteArgs) (*pb.WriteReturn, error) {
	return nfs.WriteNFS(in, wq)
}

func (s *server) Commit(ctx context.Context, in *pb.CommitArgs) (*pb.CommitReturn, error) {
	return nfs.CommitNFS(in, wq)
}

func (s *server) Readdir(ctx context.Context, in *pb.ReaddirArgs) (*pb.ReaddirReturn, error) {
	return nfs.ReaddirNFS(in)
}

func (s *server) Root(ctx context.Context, in *pb.RootArgs) (*pb.RootReturn, error) {
	return nfs.GetRootNFS(in)
}

func (s *server) Getattr(ctx context.Context, in *pb.GetAttrArgs) (*pb.GetAttrReturn, error) {
	return nfs.GetAttrNFS(in, wq)
}

func (s *server) Mkdir(ctx context.Context, in *pb.MkdirArgs) (*pb.MkdirReturn, error) {
	return nfs.MkdirNFS(in)
}

func (s *server) Rename(ctx context.Context, in *pb.RenameArgs) (*pb.RenameReturn, error) {
	return nfs.RenameNFS(in)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// check if writeverf3 hidden file exists;
	// - if it does, read from it to get number
	// - if it doesn't, write to it with the number 0
	// this will be returned on every write or commit request

	// client will read it and get in sync
	// if client has the number and it gets one back that is different, it knows
	// server crashed; it should retry everything in its buffer

	// client cleans its buffer on two conditions:
	// 1) a commit successfully returns
	// 2) the server returns back an incremented commit_number, meaning a commit happened

	// server crash number overrides commit number; we should retry everything if both
	// are returned as different

	var writeverf3 int32
	file, err := os.OpenFile(writeverf3_file_path, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("error on opening: %v\n", err)
		os.Exit(-1)
	}
	_, err = fmt.Fscanf(file, "%d", &writeverf3)
	if err != nil {
		fmt.Printf("error on scanning: %v\n", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("error on fseek: %v\n", err)
	}

	fmt.Printf("Writeverf3 value: %d\n", writeverf3)

	// increment it and persist it again
	writeverf3 += 1
	fmt.Fprintf(file, "%d", writeverf3)

	wq = nfs.InitializeServerWriteQueue(writeverf3)

	// run the server
	s := grpc.NewServer()
	pb.RegisterNFSServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
