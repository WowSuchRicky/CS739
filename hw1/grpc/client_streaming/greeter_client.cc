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

#include <time.h>
#include <grpc++/grpc++.h>

#include "helloworld.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::ClientWriter;
using grpc::Status;
using helloworld::HelloRequest;
using helloworld::HelloReply;
using helloworld::Greeter;

class GreeterClient {
 public:
  GreeterClient(std::shared_ptr<Channel> channel)
      : stub_(Greeter::NewStub(channel)) {}

  // Assambles the client's payload, sends it and presents the response back
  // from the server.
  std::string SayHello(std::string user) {
    // Data we are sending to the server.
    HelloRequest request;
    request.set_name(user);

    // Container for the data we expect from the server.
    HelloReply reply;

    // Context for the client. It could be used to convey extra information to
    // the server and/or tweak certain RPC behaviors.
    ClientContext context;

    // The actual RPC.
    Status status = stub_->SayHello(&context, request, &reply);

    // Act upon its status.
    if (status.ok()) {
      return reply.message();
    } else {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
      return "RPC failed";
    }
  }

  std::string SayHelloClientStream(std::vector<std::string> &users) {

    HelloReply reply;
    ClientContext context;

    // the object that will do the streaming for us
    std::unique_ptr<ClientWriter<HelloRequest>> writer(stub_->SayHelloClientStream(&context, &reply));

    // write to that stream
    for (int i = 0; i < users.size(); i++) {
      HelloRequest request;
      request.set_name(users[i]);
      if (!writer->Write(request)) break;
    }

    writer->WritesDone();
    Status status = writer->Finish();
    if (status.ok()) {
      //std::cout << "Finished writing\n" << "Reply is: " + reply.message() + "\n";
    } else {
      std::cout << "Failed to write\n";
    }
    return reply.message();
  }

 private:
  std::unique_ptr<Greeter::Stub> stub_;
};

int main(int argc, char** argv) {
  // Instantiate the client. It requires a channel, out of which the actual RPCs
  // are created. This channel models a connection to an endpoint (in this case,
  // localhost at port 50051). We indicate that the channel isn't authenticated
  // (use of InsecureChannelCredentials()).
  GreeterClient greeter(grpc::CreateChannel(
      "royal-02.cs.wisc.edu:50051", grpc::InsecureChannelCredentials()));
  int i, n, n_streams, n_char_per_stream; 

  if (argc < 3) {
    std::cout << "greeter_client N_STREAMS N_CHAR_PER_STREAM\n";
    return -1;
  }

  // how many characters in one of the 'streams'
  n_streams = atoi(argv[1]);

  // how many streams?
  n_char_per_stream = atoi(argv[2]);

  int total_bytes = n_streams * n_char_per_stream * 2; // factor of 2 because sending and receiving back

  std::vector<std::string> data_to_send;
  std::string to_add(n_char_per_stream, 'z');
  for (i = 0; i < n_streams; i++) {
    data_to_send.push_back(to_add);
  }

  struct timespec start, end;
  int dt;
  clock_gettime(CLOCK_REALTIME, &start);
  std::string reply = greeter.SayHelloClientStream(data_to_send);
  clock_gettime(CLOCK_REALTIME, &end);
  dt = 1000000000 * (end.tv_sec - start.tv_sec) + end.tv_nsec - start.tv_nsec;
  
  std::cout << n_streams << "," << n_char_per_stream << "," << total_bytes << "," << dt << "\n";

  return 0;
}
