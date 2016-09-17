package main
 
import (
    "fmt"
    "net"
    "os"
    "time"
)
 
/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
 
func main() {
    // Main UDP receiving port, 10001
    ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
    CheckError(err)

    // Reliability port 10002
    ReliableAddr,err := net.ResolveUDPAddr("udp",":10002")
    CheckError(err)

    // Local address
    LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
 
    /* Now listen at port 10001 */
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()

    // Open port 10002 for sending back to the client
    Conn, err := net.DialUDP("udp", LocalAddr, ReliableAddr)
    CheckError(err)
 
    buf := make([]byte, 1024)
 
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        fmt.Println("Received ",string(buf[0:n]), " from ",addr)
 
        if err != nil {
            fmt.Println("Error: ",err)
        } 

        fmt.Printf("Sending ACK...\n")

        msg := "ACK"
        buf1 := []byte(msg)

        aheh,err := Conn.Write(buf1)
        if err != nil {
            fmt.Println(msg, err)
        }
        fmt.Printf("%d\n", aheh)
        time.Sleep(time.Second * 1)


    }
}