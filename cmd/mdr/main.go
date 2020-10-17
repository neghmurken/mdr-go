package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/neghmurken/mdr/pkg/protocol"
    "github.com/neghmurken/mdr/pkg/util"
    _state "github.com/neghmurken/mdr/pkg/state"
)

func main() {
    var target string
    var port string
    var name string
    var debug bool
    
    flag.StringVar(&target, "target", "", "Host to connect to")
    flag.StringVar(&port, "port", "13000", "Listening port")
    flag.StringVar(&name, "name", "John Doe", "Your name")
    flag.BoolVar(&debug, "debug", false, "Enable debug mode")

    flag.Parse()

    ip, err := util.ExternalIP()
    if err != nil {
        util.Quit("Unable to get external IP", err)
    }

    state := _state.NewState(name, fmt.Sprintf("%s:%s", ip, port), debug)

    if (0 != len(target)) {
        Connect(target, &state)
    }

    fmt.Printf("Listening on : %s. Name : %s. Debug : %t\n", state.Authority, state.Name, state.Debug)

    go ListenForInput(func (data string) {
        state.Broadcast(data)
    })

    handle, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        util.Quit("Unable to start handle", err)
    }

    go func() { 
        for {
            conn, err := handle.Accept()
            if err != nil {
                util.Quit("Error while accepting connection", err)
            }

            RegisterPeer(conn, &state)
        }
    }()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

    <-sigc

    state.Quit()
    handle.Close()
    os.Exit(0)
}

func Connect(target string, state *_state.State) {
    conn, err := net.Dial("tcp", target)
    if err != nil {
        fmt.Println("Error while connecting to ", target, err)

        return
    }

    peer := RegisterPeer(conn, state)
    peer.SendInit(state.Name, state.Authority)
}

func RegisterPeer(conn net.Conn, state *_state.State) *_state.Peer {
    peer := _state.NewPeer(conn, "", state.Debug)
    
    if state.Debug {
        fmt.Printf("New connection with %s\n", peer.Authority())
    }

    state.Add(peer)

    go func() {
        peer.Listen(func (data string, length int) bool {
            verb, err := protocol.Parse(data)

            if err != nil {
                peer.SendErr(err.Error())
            } else {
                switch verb.Type {
                case "INIT":
                    peer.Authenticate(verb.Payload[0], verb.Payload[1])
                    fmt.Printf("%s [%s] joined the chat. %d peers\n", peer.Name, peer.Authority(), state.Count())
                    fallthrough

                case "ASK":
                    if !peer.Authenticated {
                        peer.SendErr("Not authenticated")

                        return false
                    }

                    peer.SendIdent(state.Name, state.Authority, state.AllAuthorities())

                case "IDENT":
                    peer.Authenticate(verb.Payload[0], verb.Payload[1])
                    fmt.Printf("%s [%s] joined the chat. %d peers\n", peer.Name, peer.PublicAuthority, state.Count())
                    if len(verb.Payload) > 2 {
                        for _, other := range verb.Payload[2:] {
                            if !state.Knows(other) {
                                Connect(other, state)
                            }
                        }
                    }

                case "MSG":
                    if !peer.Authenticated {
                        peer.SendErr("Not authenticated")

                        return false
                    }

                    fmt.Printf("%s > %s\n", peer.Name, verb.Payload[0])
                    peer.SendAck()

                case "QUIT":
                    fmt.Printf("%s [%s] leaved the chat. %d peers\n", peer.Name, peer.PublicAuthority, state.Count())
                    fallthrough

                case "REJECT":
                    return false

                case "ACK":                  
                    if !peer.Authenticated {
                        peer.SendErr("Not authenticated")

                        return false
                    }
                }
            }

            return true
        })

        state.Remove(peer)

        if state.Debug {
            fmt.Printf("Connection lost with %s\n", peer.Authority())
        }

        conn.Close()
    }()

    return peer
}

func ListenForInput(callback func(string)) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        callback(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("End of input. Will close")
    }
}
