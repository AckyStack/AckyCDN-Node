package main

import (
	"ackycdn-node/app"
	"ackycdn-node/app/initialization"
	"ackycdn-node/app/ssl"
	"crypto/tls"
	"github.com/anxuanzi/goutils/pkg/ftashutdown"
	"github.com/dgrr/http2"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	_ "github.com/jptosso/coraza-libinjection"
	_ "github.com/jptosso/coraza-pcre"
	"net"
)

func main() {
	slog.Info("starting ackycdn...")

	//begin initialization
	initialization.InitializeApplication()

	//start the server listening
	ListenAndServeAllNonBlock()

	ftashutdown.NewHook().Close(func() {
		slog.Info("shutting down, please wait...")
		ShutdownAll()
		slog.Info("application shutdown successfully!")
	})
}

func application() {

}

func ListenAndServeAllNonBlock() {
	go func() {
		err := app.G.FiberServer.Listen(":80")
		if err != nil {
			slog.Panic(err)
		}
	}()

	go func() {
		ln, err := net.Listen(fiber.NetworkTCP, ":443")
		if err != nil {
			slog.Panic(err)
		}
		http2.ConfigureServer(app.G.FiberServer.Server())
		err = app.G.FiberServer.Listener(tls.NewListener(ln, ssl.TlsServerConfig()))
		if err != nil {
			slog.Panic(err)
		}
	}()
}

func ShutdownAll() {
	defer app.G.MqConnection.Close()
	defer app.G.PersistenceVhostDB.Close()
	defer app.G.CacheStore.Close()
	defer app.G.FiberServer.Shutdown()
	defer app.G.SessionStorage.Close()
}
