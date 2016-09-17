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
func SendUDPMsg(conn *net.UDPConn, msg []byte) (err error){
    n,err := conn.Write(msg)

    if err != nil {
        fmt.Println(msg, err)
    }

    fmt.Printf("n:%d\n", n)

    time.Sleep(time.Second * 1)

    return err
}

//Timer/timeout function
func Timeout(timeout chan bool) {
            time.Sleep(2 * time.Second)
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

    for {
        msg :="testmsg"
        buf1 := []byte(msg)

        //Timeout channels
        ack := make(chan bool)
        timeout := make(chan bool)


        //Try to send UDP message
        SendUDPMsg(Conn, buf1);
        

        // Wait for ACK, start timer
        fmt.Printf("Waiting 2s for ACK...\n")
        go WaitForACK(ReliableConn, ack)
        go Timeout(timeout)


        for {

            select{
            case <- ack:
                    fmt.Println("ACK recv'd")
                    break
                    //ACK happened

            case <- timeout:
                    fmt.Println("Timeout, retry..")

                    SendUDPMsg(Conn, buf1);
                    go Timeout(timeout);
                    go WaitForACK(ReliableConn, ack)
                    //Timeout happened
            }
        }


    }
}