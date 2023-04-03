package main

import (
        "fmt"
        "log"
        "net"
)

type Message struct {
        from    string
        payload []byte
}

type Server struct {
        listenAddr string
        ln         net.Listener
        quitch     chan struct{}
        msgch      chan Message
}

func NewServer(listenAddr string) *Server {
        return &Server{
                listenAddr: listenAddr,
                quitch:     make(chan struct{}),
                msgch:      make(chan Message, 10),
        }
}

func (s *Server) start() error {
        ln, err := net.Listen("tcp", s.listenAddr)
        if err != nil {
                fmt.Println(err)
                return err
        }
        defer ln.Close()
        s.ln = ln

        go s.acceptLoop()

        <-s.quitch
        close(s.msgch)

        return nil

}

func (s *Server) acceptLoop() {
        for {
                conn, err := s.ln.Accept()
                if err != nil {
                        fmt.Println("accept error:", err)
                        continue
                }
                fmt.Println("new connection to the server:", conn.RemoteAddr())
                go s.readLoop(conn)

        }
}

func (s *Server) readLoop(con net.Conn) {
        defer con.Close()
        buf := make([]byte, 2048)

        for {
                n, err := con.Read(buf)
                if err != nil {
                        fmt.Println("read error:", err)
                        continue
                }
                s.msgch <- Message{
                        from:    con.RemoteAddr().String(),
                        payload: buf[:n],
                }

                con.Write([]byte("thank you for your message!\n"))
        }
}

func main() {

        server := NewServer(":9090")
        go func() {
                for msg := range server.msgch {
                        fmt.Printf("received message from connection (%s): %s\n", msg.from, string(msg.payload))
                }
        }()

        log.Fatal(server.start())

}