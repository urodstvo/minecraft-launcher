package launcher

import (
	_ "embed"
	"fmt"

	"github.com/urodstvo/minecraft-launcher/utils"
	webview "github.com/webview/webview_go"
)

type UIView struct {
	Logger utils.Logger

	WV webview.WebView
}

func (ui *UIView) Open(url string) {
	ui.Bind()

	ui.WV.SetTitle("Minecraft launcher")
	ui.WV.SetSize(800, 600, webview.HintFixed)

	ui.WV.Navigate(url)
	ui.WV.Run()
	defer ui.WV.Destroy()
}

func (ui *UIView) Bind() {
	ui.WV.Bind("initLauncher", func() string {		
		defer ui.Logger.Info("Launcher is initialised")
		return "0.0.0-alpha"
	})

	ui.WV.Bind("launchMinecraft", func() {
		ui.Logger.Info("Minecraft is Launched")
		ui.WV.Dispatch(func() {
			ui.WV.Terminate()
		})
	})
}

func (ui *UIView) SetVersion(version string) {
	ui.Logger.Info("Changed version")
	ui.WV.Eval(fmt.Sprintf("ui.setVersion(%s)", version))
}



