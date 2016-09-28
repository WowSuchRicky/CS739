These measurements were performed with UDP payload of 512-bytes
- 512-byte UDP payload can typically let us avoid datagram fragmentation
- see http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet

The different machines experiment were performed using betelgeuse.cs.wisc.edu as the server (receiver), and royal-06.cs.wisc.edu as the client (sender).
Compiler optimizations were disabled.
