#include <iostream>
#include <string>

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

int main(int argc, char** argv) {
  boost::shared_ptr<TTransport> socket(new TSocket("localhost", 9090));
  boost::shared_ptr<TTransport> transport(new TBufferedTransport(socket));
  boost::shared_ptr<TProtocol> protocol(new TBinaryProtocol(transport));
  CalculatorClient client(protocol);

  try {
    transport->open();

    int i, mode, n;
    std::string::size_type sz;
    mode = atoi(argv[1]);
    n = atoi(argv[2]);

    if (mode == 0) {
        int num = atoi(argv[3]);
        for (i = 0; i < n; i++)  {
            cout << "int time = " << client.int_time(num) << endl;
        }
    } else if (mode == 1) {
        double d = std::stod(argv[3], &sz);
        for (i = 0; i < n; i++) {
            cout << "dbl time = " << client.dbl_time(d) << endl;
        }
    } else if (mode == 2) {
        int l = atoi(argv[3]);
        std::string str(l, 'a');
        for (i = 0; i < n; i++) {
            cout << "str time = " << client.str_time(str) << endl;
        }
    } else {
        Blah b = Blah();
        b.num1 = atoi(argv[3]);
        b.num2 = std::stod(argv[4], &sz);
        cout << "bla time = " << client.bla_time(b) << endl;
    }
    transport->close();
  } catch (TException& tx) {
    cout << "ERROR: " << tx.what() << endl;
  }
}
