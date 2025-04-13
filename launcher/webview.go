package launcher

import (
	webview "github.com/webview/webview_go"
)

type UIView struct {
	WV webview.WebView
}

func (v *UIView) Open() {
	v.WV.Run()
}