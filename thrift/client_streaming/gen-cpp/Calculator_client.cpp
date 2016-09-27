#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

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
    long int num_nums = atof(argv[2]); // How many numbers to send
    //cout << "Summing:  " << endl;
    vector<long int> nums_vec(num_nums);
    for (long int i = 0; i < num_nums; i++){
      nums_vec.push_back(i);
    }


    /* Begin time */
    int start_time_s, total_time_s;
    unsigned long long start_time_ns, total_time_ns;
    struct timespec time;
    assert(clock_getres(CLOCK_REALTIME, &time) != -1);
    assert(clock_gettime(CLOCK_REALTIME, &time) != -1);
    
    start_time_ns = time.tv_nsec;
    start_time_s = time.tv_sec;
    client.sum(nums_vec);  

    assert(clock_gettime(CLOCK_REALTIME, &time) != -1);
    total_time_ns = time.tv_nsec - start_time_ns;
    total_time_s = time.tv_sec - start_time_s;
    
    printf("sent: %lu\n", nums_vec.size() * sizeof(long int) );
    printf("%d sec, %llu ns\n", total_time_s, total_time_ns);


    /* End time */ 

    transport->close();
  } catch (TException& tx) {
    cout << "ERROR: " << tx.what() << endl;
  }
}
