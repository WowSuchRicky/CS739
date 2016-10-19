/*
 *
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
	nfs "github.com/Ricky54326/CS739/hw2/server/lib"
	"log"
	"net"

	pb "github.com/Ricky54326/CS739/hw2/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

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
	return nfs.ReadNFS(in)
}

func (s *server) Write(ctx context.Context, in *pb.WriteArgs) (*pb.WriteReturn, error) {
	return nfs.WriteNFS(in)
}

func main() {

	//path, _ := nfs.InumToPath(663418)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNFSServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
