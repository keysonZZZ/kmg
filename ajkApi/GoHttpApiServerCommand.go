package ajkApi

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/kmgCrypto"
	"net"
	"net/http"
)

var AdditionHttpHandler []HttpHandlerConfig

type HttpHandlerConfig struct {
	Path    string
	Handler http.Handler
}

//start a golang http api server
type GoHttpApiServerCommand struct {
	Container     *dependencyInjection.Container
	http          string
	https         string
	randPort      bool
	isHttps       bool
	tcpListenAddr string
}

func (command *GoHttpApiServerCommand) SetContainer(Container *dependencyInjection.Container) {
	command.Container = Container
}
func (command *GoHttpApiServerCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoHttpApiServer", Short: `start a golang http api server `}
}
func (command *GoHttpApiServerCommand) ConfigFlagSet(f *flag.FlagSet) {
	f.StringVar(&command.http, "http", ":18080", "listen http port of the server")
	f.StringVar(&command.https, "https", "", "listen https port of the server")
	f.BoolVar(&command.randPort, "randPort", false, "if can not listen on default port ,will listen on random port")
}

func (command *GoHttpApiServerCommand) Execute(context *console.Context) error {
	if command.https != "" {
		command.isHttps = true
		command.tcpListenAddr = command.https
	} else {
		command.tcpListenAddr = command.http
	}
	c := command.Container
	ihandler, err := c.GetByType((*JsonHttpHandler)(nil))
	if err != nil {
		return err
	}
	jsonHttpHandler := ihandler.(*JsonHttpHandler)
	http.Handle("/api", &HttpApiFilterManager{
		Filters: []HttpApiFilter{
			jsonHttpHandler.Filter,
		},
	})
	http.Handle("/api.deflate", &HttpApiFilterManager{
		Filters: []HttpApiFilter{
			HttpApiDeflateCompressFilter,
			jsonHttpHandler.Filter,
		},
	})
	for _, handlerConfig := range AdditionHttpHandler {
		http.Handle(handlerConfig.Path, handlerConfig.Handler)
	}
	l, err := command.listen()
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, "Listen on %s\n", l.Addr().String())
	if command.isHttps {
		tlsConfig, err := kmgCrypto.CreateTlsConfig()
		if err != nil {
			return fmt.Errorf("fail at kmgTls.CreateTlsConfig,error:%s", err.Error())
		}
		l = tls.NewListener(l, tlsConfig)
	}
	return http.Serve(l, nil)
}

//first try addr,if err happened try random address.
func (command *GoHttpApiServerCommand) listen() (l net.Listener, err error) {
	l, err = net.Listen("tcp", command.tcpListenAddr)
	if err == nil {
		return
	}
	if command.randPort {
		l, err = net.Listen("tcp", ":0")
		return
	}
	return
}
