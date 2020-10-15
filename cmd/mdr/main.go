package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "os"
)

func main() {
    var target string
    var port string
    
    flag.StringVar(&target, "target", "", "Host to connect to")
    flag.StringVar(&port, "port", "13000", "Listening port")
    flag.Parse()

    if (0 != len(target)) {
        go connect(target)
    }

    fmt.Println("Listening on :", port)

    handle, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        quit("Unable to start handle", err)
    }

    defer handle.Close()

    for {
        conn, err := handle.Accept()
        if err != nil {
            quit("Error while accepting connection", err)
        }

        go handleConnection(conn)
    }
}

func connect(target string) {
    conn, err := net.Dial("tcp", target)
    if err != nil {
        fmt.Println("Error while connecting to ", target, err)

        return
    }

    defer conn.Close()

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        fmt.Fprintf(conn, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("End of input. Will close")
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

    var buff = make([]byte, 1024)

    for {
        byteRead, err := conn.Read(buff)

        if err != nil {
            fmt.Println("Error while reading connection.", err)

            break
        }

        fmt.Printf("Received [%d] : %s\n", byteRead, string(buff))
    }
}

func quit(message string, err error) {
    fmt.Println(message)
    if nil != err {
        fmt.Println(err)
    }

    flag.PrintDefaults()
    os.Exit(1)
}
