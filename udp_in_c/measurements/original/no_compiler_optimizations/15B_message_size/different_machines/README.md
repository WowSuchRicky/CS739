All of the measurements here were taken using:
- a 1ms timeout period (the period between send attempts)
- royal-05.cs.wisc.edu as the client (sender) - kernel: Linux royal-05 3.13.0-95-generic #142-Ubuntu (UBUNTU)
- betelgeuse.cs.wisc.edu as the server (receiver) - kernel: Linux betelgeuse.cs.wisc.edu 2.6.32-573.7.1.el6.x86_64 (REDHAT)
- the message being sent is 15 bytes
- the acknowledgment is a single byte
- no optimizations (passing -O0 into gcc)
