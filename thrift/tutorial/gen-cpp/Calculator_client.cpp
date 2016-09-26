#include <iostream>
#include <stdio.h>
#include <stdlib.h>

#include <thrift/protocol/TBinaryProtocol.h>
#include <thrift/transport/TSocket.h>
#include <thrift/transport/TTransportUtils.h>

#include "../gen-cpp/Calculator.h"

using namespace std;
using namespace apache::thrift;
using namespace apache::thrift::protocol;
using namespace apache::thrift::transport;

using namespace tutorial;
using namespace shared;

int main(int argc, char *argv[]) {
  std::string server = argv[1];

  boost::shared_ptr<TTransport> socket(new TSocket(server, 9090));
  boost::shared_ptr<TTransport> transport(new TBufferedTransport(socket));
  boost::shared_ptr<TProtocol> protocol(new TBinaryProtocol(transport));
  CalculatorClient client(protocol);

  try {
    transport->open();



    /* Begin sum() call on the server */
    double num_nums = atof(argv[2]); // How many numbers to send
    cout << "Summing:  " << endl;
    vector<long int> nums_vec(num_nums);
    for (double i = 0; i < num_nums; i++){
      nums_vec.push_back(i);
    }

    client.sum(nums_vec);
    

    transport->close();
  } catch (TException& tx) {
    cout << "ERROR: " << tx.what() << endl;
  }
}
