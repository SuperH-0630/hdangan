package v1main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"github.com/SuperH-0630/hdangan/src/systeminit"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	runtimeos "runtime"
	"time"
)

func NewInitFile(rt runtime.RunTime, w fyne.Window) {
	c, err := systeminit.GetInit()
	if err != nil {
		dialog.ShowError(fmt.Errorf("配置文件系统错误: %s", err.Error()), w)
		return
	}

	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
			return
		}

		defer func() {
			_ = reader.Close()
		}()

		fileByte, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件读取有误，不能使用: %s", err.Error()), w)
			return
		}

		fileString := string(fileByte)

		err = systeminit.CheckInit(fileString)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件有误，不能使用: %s", err.Error()), w)
			return
		}

		dialog.ShowConfirm("确认？", "是否确认更换配置文件？", func(b bool) {
			if !b {
				return
			}

			if systeminit.IsPathExists(c.ConfigDir) {
				oldConfig, err := os.ReadFile(c.ConfigDir)
				if err != nil {
					dialog.ShowError(fmt.Errorf("无法正确备份旧配置: %s", err.Error()), w)
					return
				}

				newPath := c.ConfigDir + time.Now().Format("20060102150405")
				err = os.WriteFile(newPath, oldConfig, 0666)
				if err != nil {
					dialog.ShowError(fmt.Errorf("无法正确备份旧配置: %s", err.Error()), w)
					return
				}

				err = os.Remove(c.ConfigDir)
				if err != nil {
					dialog.ShowError(fmt.Errorf("无法正确备份旧配置: %s", err.Error()), w)
					return
				}
			}

			err := os.WriteFile(c.ConfigDir, fileByte, 0666)
			if err != nil {
				dialog.ShowError(fmt.Errorf("新配置文件写入可能产生错误: %s", err.Error()), w)
				return
			}

			err = systeminit.ReInit()
			if err != nil {
				dialog.ShowError(fmt.Errorf("配置文件已经呗替换，但是无法通过系统测试: %s", err.Error()), w)
				return
			}

			dialog.ShowInformation("提示", "新的配置文件已经导入并重载，但请尽量重新打开软件也防止未知问题的发生。", w)
		}, w)

	}, w)
}

func OpenInit(rt runtime.RunTime, w fyne.Window) {
	c, err := systeminit.GetInit()
	if err != nil {
		dialog.ShowError(fmt.Errorf("配置文件系统错误: %s", err.Error()), w)
		return
	}

	if !systeminit.IsPathExists(c.ConfigDir) {
		dialog.ShowError(fmt.Errorf("系统配置文件不存在"), w)
		return
	}

	var cmd *exec.Cmd
	sysType := runtimeos.GOOS
	if sysType == "windows" {
		cmd = exec.Command("cmd", "/c", "start", c.ConfigDir)
	} else if sysType == "darwin" {
		cmd = exec.Command("open", c.ConfigDir)
	} else {
		cmd = exec.Command("xdg-open", c.ConfigDir)
	}

	err = cmd.Start()
	if err != nil {
		dialog.ShowError(fmt.Errorf("操作系统不支持"), w)
	}
}

func SaveInit(rt runtime.RunTime, w fyne.Window) {
	c, err := systeminit.GetInit()
	if err != nil {
		dialog.ShowError(fmt.Errorf("配置文件系统错误: %s", err.Error()), w)
		return
	}

	if !systeminit.IsPathExists(c.ConfigDir) {
		dialog.ShowError(fmt.Errorf("系统配置文件不存在"), w)
		return
	}

	dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
			return
		}

		defer func() {
			_ = writer.Close()
		}()

		data, err := yaml.Marshal(c.Yaml)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件转义错误: %s", err.Error()), w)
		}

		_, err = writer.Write(data)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件保存错误: %s", err.Error()), w)
		}

		return
	}, w)
	dlg.SetFileName("config.yaml")
	dlg.Show()
}

func CopyInit(rt runtime.RunTime, w fyne.Window) {
	c, err := systeminit.GetInit()
	if err != nil {
		dialog.ShowError(fmt.Errorf("配置文件系统错误: %s", err.Error()), w)
		return
	}

	if !systeminit.IsPathExists(c.ConfigDir) {
		dialog.ShowError(fmt.Errorf("系统配置文件不存在"), w)
		return
	}

	dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer == nil {
			return
		} else if err != nil {
			dialog.ShowError(fmt.Errorf("选择框遇到错误：%s", err.Error()), w)
			return
		}

		defer func() {
			_ = writer.Close()
		}()

		data, err := os.ReadFile(c.ConfigDir)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件读取错误: %s", err.Error()), w)
		}

		_, err = writer.Write(data)
		if err != nil {
			dialog.ShowError(fmt.Errorf("配置文件保存错误: %s", err.Error()), w)
		}

		return
	}, w)

	dlg.SetFileName("config-copy.yaml")
	dlg.Show()
}
