package main
 
import (
    "fmt"
    "net"
    "time"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

// Wait for ACK
func WaitForACK(conn *net.UDPConn, ack chan bool){

    //Allocate 100 bytes for recv
    buf := make([]byte, 100)

    n,addr,_ := conn.ReadFromUDP(buf)
    fmt.Printf("nbytes: %d addr: %s\n", n, addr)
    ack <- true
}

//Send UDP msg
func SendUDPMsg(conn *net.UDPConn, ReliableConn *net.UDPConn, msg []byte) (err error){
    for {


        n,err := conn.Write(msg)

        //Timeout channels
        ack := make(chan bool)
        timeout := make(chan bool)

        if err != nil {
            fmt.Println(msg, err)
        }

        fmt.Printf("written:%d\n", n)


     // Wait for ACK, start timer
        fmt.Printf("Waiting 5s for ACK...\n")
        go WaitForACK(ReliableConn, ack)
        go Timeout(timeout)


        
            select{
            case <- ack:
                    fmt.Println("ACK recv'd\n")
                    return nil
                    //ACK happened

            case <- timeout:
                    fmt.Println("Timeout, retry..\n")
                    //Timeout happened
            }
        
        }

    return err
}

//Timer/timeout function
func Timeout(timeout chan bool) {
            time.Sleep(5 * time.Second)
            timeout <- true
}


func main() {

    //Store local address for sending, regular UDP comms
    ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:10001")
    CheckError(err)

    //Open another UDP socket for sending, reliable mechanisms
    ServerAddr2,err := net.ResolveUDPAddr("udp","127.0.0.1:10002")
    CheckError(err)
 
    //Store local addr for listening
    LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
 
    //"Dial" server
    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)

    /* Now listen at selected port */
    ReliableConn, err := net.ListenUDP("udp", ServerAddr2)
    CheckError(err)
    defer ReliableConn.Close()
 
    defer Conn.Close()

    // Infinitely send msgs, for now
    i := 0
    for {
        i++
        msg := fmt.Sprintf("Test: #%d\n", i)
        buf1 := []byte(msg)

        //Try to send UDP message
        SendUDPMsg(Conn, ReliableConn, buf1);

    }
}