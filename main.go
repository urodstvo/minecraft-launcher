package main

import (
	"github.com/urodstvo/minecraft-launcher/launcher"
	webview "github.com/webview/webview_go"
)

func main(){
	ui := &launcher.UIView{
		WV: webview.New(false),
	}
	l := launcher.Launcher{
		UI: ui,
	}
	l.Run()
}