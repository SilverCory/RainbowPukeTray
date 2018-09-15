package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/SilverCory/RainbowPukeTray"
	"github.com/getlantern/systray"
	"github.com/shibukawa/configdir"
)

type Config struct {
	On bool `json:"on"`
}

const TextTurnOn = "Turn On"
const TextTurnOff = "Turn Off"

var cmd *exec.Cmd

func main() {

	// config setup.
	config := Config{true}
	configDirs := configdir.New("rainbowpuke", "rainbowpuketray")
	configDirs.LocalPath, _ = filepath.Abs(".")
	folder := configDirs.QueryFolderContainsFile("setting.json")
	if folder != nil {
		data, _ := folder.ReadFile("setting.json")
		json.Unmarshal(data, &config)
	} else {
		data, _ := json.Marshal(&config)
		folders := configDirs.QueryFolders(configdir.Global)
		folders[0].WriteFile("setting.json", data)
	}

	systray.Run(func() {
		fmt.Println("READY")
		var toggleItem = systray.AddMenuItem("aaaaa", "Change the state of the rainbow puke.")
		cmd = exec.Command("RainbowPuke.exe")

		if config.On {
			turnOn(toggleItem)
		} else {
			turnOff(toggleItem)
		}

		go screenChangeRestart(&config, toggleItem)

		go func() {
			for {
				select {
				case <-toggleItem.ClickedCh:
					config.On = !config.On
					if config.On {
						turnOn(toggleItem)
					} else {
						turnOff(toggleItem)
					}
				}
			}
		}()

		systray.AddSeparator()
		systray.AddSeparator()
		systray.AddSeparator()

		exitItem := systray.AddMenuItem("Exit", "Off I fuck...")
		go func() {
			<-exitItem.ClickedCh
			systray.Quit()
			turnOff(toggleItem)
		}()
	}, func() {
		fmt.Println("Cya!")
		data, _ := json.Marshal(&config)
		folders := configDirs.QueryFolders(configdir.Global)
		folders[0].WriteFile("setting.json", data)
	})

}

func screenChangeRestart(config *Config, toggleItem *systray.MenuItem) {

	var (
		user32           = syscall.NewLazyDLL("User32.dll")
		getSystemMetrics = user32.NewProc("GetSystemMetrics")
	)

	const (
		SM_CXVIRTUALSCREEN uintptr = 78
		SM_CYVIRTUALSCREEN         = 79
	)

	lastX, lastY := 0, 0
	for {

		x, _, err := getSystemMetrics.Call(SM_CXVIRTUALSCREEN)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}

		y, _, err := getSystemMetrics.Call(SM_CYVIRTUALSCREEN)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}

		if lastX == 0 {
			lastX = int(x)
		}
		if lastY == 0 {
			lastY = int(y)
		}

		if lastX != int(x) || lastY != int(y) {
			fmt.Println("SCREEN SIZE CHANGE FOUND!")
			if config.On {
				turnOff(toggleItem)
				time.Sleep(100 * time.Millisecond)
				turnOn(toggleItem)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func turnOff(toggleItem *systray.MenuItem) {
	if cmd.Process != nil {
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			fmt.Println(err)
			cmd.Process.Kill()
		}
	}
	toggleItem.SetTitle(TextTurnOn)
	systray.SetIcon(RainbowPukeTray.GreyIcon)
}

func turnOn(toggleItem *systray.MenuItem) {
	cmd = exec.Command("RainbowPuke.exe")
	startCommand(cmd, func(e error) {
		turnOff(toggleItem)
	})
	toggleItem.SetTitle(TextTurnOff)
	systray.SetIcon(RainbowPukeTray.Icon)
}

func startCommand(cmd *exec.Cmd, done func(error)) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	go func() {
		done(cmd.Run())
	}()
}
