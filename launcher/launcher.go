package launcher

type Launcher struct {
	Version string
	System  *System

	Server    *Server
	Minecraft *Minecraft
	UI        *UIView
}

func (l *Launcher) checkVersion() {}

func (l *Launcher) Run() {
	l.UI.Open()
}