package launcher

import (
	"github.com/urodstvo/minecraft-launcher/utils"
)


type Launcher struct {
	Logger utils.Logger

	Version string
	System  *System

	Server 	  *Server
	Minecraft *Minecraft
	UI        *UIView
}

func (l *Launcher) checkVersion() {
	// lastVersion := "0.0.0-beta"
	
}

func (l *Launcher) Run() {
	l.Logger.Info("Launcher started")

	l.checkVersion()
	url := l.Server.Run()
	l.UI.Open(url)
	
	l.Logger.Info("Launcher stoped")
}