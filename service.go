package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"github.com/shazow/go-irckit"
	"github.com/thoj/go-ircevent"
	irct "github.com/sorcix/irc"
)

// version gets replaced during build
var version string = "dev"
var federated map[string]*irc.Connection
// logger gets replaced by golog

// Options contains the flag options
type Options struct {
	Bind    string `long:"bind" description:"Bind address to listen on." value-name:"[HOST]:PORT" default:":6667"`
	Pprof   string `long:"pprof" description:"Bind address to serve pprof for profiling." value-name:"[HOST]:PORT"`
	Name    string `long:"name" description:"Server name." default:"irckit-demo"`
	Motd    string `long:"motd" description:"Message of the day."`
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Version bool   `long:"version"`
}


func fail(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}
func newChannel2(server irckit.Server, name string) irckit.Channel {
	ret := irckit.NewChannel(server, name)
	parts := strings.Split(name, "@")
	fmt.Printf("New Channel: %s\n", name)
	if len(parts) > 1 {
		fmt.Printf("New federated channel on: %s\n", parts[1])
		go func() { if federated[parts[1]] == nil {
			q := irc.IRC(strings.Replace(server.Name(), ".", "-", -1), "federated")
			federated[parts[1]] = q
			federated[parts[1]].AddCallback("001", func(e *irc.Event) {
				federated[parts[1]].Join(parts[0])
				
			})
			federated[parts[1]].AddCallback("PRIVMSG", func (e *irc.Event) {
				var fakeuser irckit.User
				fakeuser.Nick = e.Nick
				fakeuser.User = e.Nick
				fakeuser.Real = e.Nick
				var themsg string;
				if e.User != "federated" {
				fakeuser.Host = strings.Replace(parts[1], "!", "-", -1)
				themsg = e.Message()
				} else {
				themsgparts := strings.Split(e.Message(), ": ")
				themsg = themsgparts[1]
				fakeuser.Nick = themsgparts[0]
				fakeuser.Host = strings.Replace(e.Nick, "-", ".", -1)
				}
				var sndmsg irct.Message;
				sndmsg.Trailing = themsg;
				fmt.Println(e.Arguments)
				sndmsg.Params = []string{e.Arguments[0]}
				if toChan, exists := server.HasChannel(e.Arguments[0] + "@" + parts[1]); exists {
				fmt.Println("Sending");
				fmt.Println(e.Arguments[0])
		toChan.Message(&fakeuser, themsg)
		
			}
			})
	
			
			federated[parts[1]].Connect(strings.Replace(parts[1],"!",":",-1))
			federated[parts[1]].Loop()
			
		} else {

				federated[parts[1]].Join(parts[0])
				

		}

		}()
	}
	return ret
}
func MitmPrivMsg(s irckit.Server, u *irckit.User, msg *irct.Message) error {
	om := msg.Trailing
	fmt.Println(om)
	parts := strings.Split(msg.Params[0], "@")
	if parts[0][0] == '#' {
		if parts[1] != "" {
			fmt.Println("Remote Send")
			if federated[parts[1]] != nil {
		fmt.Println(om)
federated[parts[1]].Privmsgf(parts[0], "%s: %s", u.Nick, om)
			}
		}
	}
	return irckit.CmdPrivMsg(s, u, msg)
}

	type cmdsx map[string]irckit.Handler
func main() {
	federated = make(map[string]*irc.Connection)
	options := Options{}
	options.Bind = ":" + os.Args[3]
	options.Motd = "Hello World"
	options.Name = "server.net"




	if options.Pprof != "" {
		go func() {
			fmt.Println(http.ListenAndServe(options.Pprof, nil))
		}()
	}

	// Figure out the log level



	
	cmds := irckit.DefaultCommands()
	
	cmds.Add( irckit.Handler{Command: irct.PRIVMSG, Call: MitmPrivMsg, MinParams: 1})


	socket, err := net.Listen("tcp", options.Bind)
	if err != nil {
		fail(4, "Failed to listen on socket: %v\n", err)
	}
	defer socket.Close()

	motd := []string{}
	if options.Motd != "" {
		motd = append(motd, options.Motd)
	}
	srv := irckit.ServerConfig{
		Name: os.Args[1],
		Motd: []string{os.Args[2], "+federated=1"},
		Version: "1.0",
		NewChannel: newChannel2,
		Commands: cmds,
	}.Server()
	go start(srv, socket)

	fmt.Printf("Listening for connections on %v\n", socket.Addr().String())

	// Construct interrupt handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig // Wait for ^C signal
	fmt.Fprintln(os.Stderr, "Interrupt signal detected, shutting down.")
	srv.Close()
	os.Exit(0)
}

func start(srv irckit.Server, socket net.Listener) {
	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			return
		}

		// Goroutineify to resume accepting sockets early
		go func() {
			fmt.Printf("New connection: %s\n", conn.RemoteAddr())
			err = srv.Connect(irckit.NewUserNet(conn))
			if err != nil {
				fmt.Printf("Failed to join: %v\n", err)
				return
			}
		}()
	}
}
