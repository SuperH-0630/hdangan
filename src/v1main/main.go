package v1main

import (
	"fmt"
	happ "github.com/SuperH-0630/hdangan/src/app"
	"github.com/SuperH-0630/hdangan/src/assest"
	"github.com/SuperH-0630/hdangan/src/fail"
	"github.com/SuperH-0630/hdangan/src/model"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/flopp/go-findfont"
	"os"
	"strings"
)

var FONT_NAME = ""

// ALLOW_FONT_NAME 常见字体：宋体（SimSun, STSong）、（SimHei, STHei）、微软雅黑（Microsoft YaHei）、Linux的其他字体
var ALLOW_FONT_NAME = []string{"Dengb.ttf", "Deng.ttf", "Dengl.ttf", "KaiTi.ttf", "SimSun.ttc", "STSong.ttf", "SimHei.ttc", "STHei.ttf", "Microsoft YaHei",
	"WenQuanYi", "Source", "Fangzheng", "Arphic", "AR PL", "ZCOOL"}

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		for _, n := range ALLOW_FONT_NAME {
			if strings.Contains(strings.ToLower(path), strings.ToLower(n)) {
				FONT_NAME = path
				return
			}
		}
	}
	fail.ToFail("Font not found.")
}

func Main() {
	start()

	a := happ.NewApp()
	a.SetIcon(assest.MainIco)

	rt := runtime.NewRunTime(a)

	err := model.AutoCreateModel(rt)
	if err != nil {
		fail.ToFail(fmt.Sprintf("数据库构建失败: %s。", err.Error()))
		return
	}

	ShowHelloWindow(rt)
	StartTheGame(rt)

	a.Run()

	exit()
}

func start() {
	err := os.Setenv("FYNE_FONT", FONT_NAME)
	if err != nil {
		fail.ToFail("font setting failed. 字体设置失败。")
		return
	}
}

func exit() {
	_ = os.Unsetenv("FYNE_FONT")
	fmt.Print("EXIT")
}
