
all: server client

server:
	g++ -o Calculator_server -I../gen-cpp -I/unsup/boost-1.58.0/include -I/unsup/thrift-0.9.3/include shared_types.cpp SharedService.cpp tutorial_types.cpp Calculator.cpp Calculator_server.cpp -L/unsup/boost-1.58.0/lib -L/unsup/thrift-0.9.3/amd64_ubu14/lib/ -lthrift

client:
	g++ -o Calculator_client -I../gen-cpp -I/unsup/boost-1.58.0/include -I/unsup/thrift-0.9.3/include shared_types.cpp SharedService.cpp tutorial_types.cpp Calculator.cpp Calculator_client.cpp -L/unsup/boost-1.58.0/lib -L/unsup/thrift-0.9.3/amd64_ubu14/lib/ -lthrift

clean:
	 rm -f Calculator_server Calculator_client
