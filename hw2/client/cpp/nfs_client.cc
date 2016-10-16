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

#include <iostream>
#include <memory>
#include <string>

#include <grpc++/grpc++.h>

#include "nfs.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using nfs::FileHandle;
using nfs::Attribute;
using nfs::LookupArgs;
using nfs::LookupReturn;
using nfs::CreateArgs;
using nfs::CreateReturn;
using nfs::RemoveArgs;
using nfs::RemoveReturn;
using nfs::ReadArgs;
using nfs::ReadReturn;
using nfs::WriteArgs;
using nfs::WriteReturn;
using nfs::NFS;

class NFSClient {
 public:
  NFSClient(std::shared_ptr<Channel> channel)
      : stub_(NFS::NewStub(channel)) {}

  // Assembles the client's payload, sends it and presents the response back
  // from the server.
  LookupReturn lookup(FileHandle& dirfh, const std::string& name) {
    std::cout << "In lookup\n";
    // Data we are sending to the server.
    LookupArgs args;
    std::cout << "dirfh.inode(): " << dirfh.inode() << std::endl;
    std::cout << "dirfh.fsnum(): " << dirfh.fsnum() << std::endl;
    std::cout << "dirfh.genum(): " << dirfh.genum() << std::endl;
    args.set_allocated_dirfh(&dirfh);
    args.set_name(name);

    // Container for the data we expect from the server.
    LookupReturn reply;

    // Context for the client. It could be used to convey extra information to
    // the server and/or tweak certain RPC behaviors.
    ClientContext context;

    // The actual RPC.
    Status status = stub_->lookup(&context, args, &reply);

    std::cout << "In lookup, stub returned\n";
    // Act upon its status.
    if (status.ok()) {
      std::cout << "Success!\n";
      return reply;
    } else {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
      return reply; //TODO: Maybe return some actual error/exception?
    }
  }

 private:
  std::unique_ptr<NFS::Stub> stub_;
};

int main(int argc, char** argv) {
  // Instantiate the client. It requires a channel, out of which the actual RPCs
  // are created. This channel models a connection to an endpoint (in this case,
  // localhost at port 50053). We indicate that the channel isn't authenticated
  // (use of InsecureChannelCredentials()).
  NFSClient greeter(grpc::CreateChannel(
      "localhost:50053", grpc::InsecureChannelCredentials()));
  std::string name("world");
  FileHandle dirfh;
  dirfh.set_inode(1);
  dirfh.set_fsnum(2);
  dirfh.set_genum(3);
  LookupReturn reply = greeter.lookup(dirfh, name);

  return 0;
}
