package initialization

import (
	"ackycdn-node/app"
	"ackycdn-node/app/types"
	"ackycdn-node/app/vhost"
	"errors"
	"github.com/asdine/storm/v3"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/slog"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func InitializeApplication() {
	slog.Info("starting application...")
	app.G = &app.GlobalResource{
		FiberServer:        nil,
		CdnCache:           nil,
		PersistenceVhostDB: nil,
		VhostConfigsMem:    nil,
		MqConnection:       nil,
		Mq:                 nil,
		NodeInfo:           nil,
	}
	//initialization OS
	slog.Info("init  OS environment...")
	slog.SetLogLevel(slog.DebugLevel)
	if strings.Contains(runtime.GOOS, "linux") {
		err := initOSEnv()
		if err != nil {
			slog.Panic(err)
		}
		slog.Info("init OS environment... done")
	} else {
		slog.Info("init  OS environment... skipped")
	}

	slog.Info("init  storage...")
	initStores()
	slog.Info("init  storage... done")

	slog.Info("init  configurations...")
	initConfigurations()
	slog.Info("init  configurations... done")

	slog.Info("init  MQs...")
	initMq()
	slog.Info("init  MQs... done")

	slog.Info("init  WAF features...")
	initWaf()
	slog.Info("init  WAF features... done")

	slog.Info("init  server...")
	initFiberServer()
	slog.Info("init  server... done")
}

func initConfigurations() {
	//check and register node if needed
	nodeInfo := &types.NodeConfig{}
	err := app.G.PersistenceVhostDB.One("CfgKey", "MainNodeConfig", nodeInfo)
	if err != nil {
		if err == storm.ErrNotFound {
			//TODO
			//request for a new ENID
			nodeInfo.CfgKey = "MainNodeConfig"
			nodeInfo.CreateTime = carbon.Now().TimestampWithMillisecond()
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				slog.Panic(err)
			}
			defer conn.Close()
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			nodeInfo.MainIP = localAddr.IP.String()
			nodeInfo.NodeId = "07KG19WTNMJH" //12 chars

			err = app.G.PersistenceVhostDB.Save(nodeInfo)
			if err != nil {
				slog.Panic(err)
			}
		} else {
			slog.Panic(err)
		}
	}
	app.G.NodeInfo = nodeInfo

	//save vhost cfg from db to mem
	var vhosts []*types.VHostConfig
	err = app.G.PersistenceVhostDB.All(&vhosts)
	if err != nil {
		slog.Panic(err)
	}
	for _, vh := range vhosts {
		vhost.PutConfigMem(vh)
	}
}

// initOSEnv set environments
func initOSEnv() error {
	rLimit := syscall.Rlimit{Cur: 1024000, Max: 1024000}
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	cmd := exec.Command("/usr/sbin/sysctl", "-w", "net.core.somaxconn=65535")
	err = cmd.Run()
	if err != nil {
		cmd = exec.Command("/sbin/sysctl", "-w", "net.core.somaxconn=65535")
		err = cmd.Run()
		if err != nil {
			return errors.New("sysctl set net.core.somaxconn error:" + err.Error())
		}
	}
	cmd = exec.Command("/usr/sbin/sysctl", "-w", "net.ipv4.tcp_max_syn_backlog=1024000")
	err = cmd.Run()
	if err != nil {
		cmd = exec.Command("/sbin/sysctl", "-w", "net.ipv4.tcp_max_syn_backlog=1024000")
		err = cmd.Run()
		if err != nil {
			return errors.New("sysctl set net.ipv4.tcp_max_syn_backlog error: " + err.Error())
		}
	}
	return nil
}
