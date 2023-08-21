package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/getlantern/systray"
	"gopkg.in/yaml.v2"
)

var fileConfigPath string

type RootConfig struct {
	Root []FileConfig `yaml:"root"`
}

type FileConfig struct {
	Label    string       `yaml:"label"`
	Run      string       `yaml:"run"`
	Divider  bool         `yaml:"divider"`
	Children []FileConfig `yaml:"children"`
}

func main() {
	fileConfigPath = os.Getenv("HOME") + "/sh-icon-tray.yml"
	CreateFileConfig()
	systray.Run(onReady, func() {})
}

func onReady() {
	rootConfig := ReadFileConfig()
	systray.SetIcon(IconTerminal)
	systray.SetTitle("Bash")
	systray.SetTooltip("Run your scripts :)")
	CreateMenuRecursive(rootConfig.Root, nil)
	systray.AddSeparator()
	moreDropDown := systray.AddMenuItem("More", "")
	moreEditConfig := moreDropDown.AddSubMenuItem("Edit Config", "")
	go RunScript(moreEditConfig, FileConfig{
		Label: "Edit Config",
		Run:   "code " + fileConfigPath,
	})

	moreRefresh := moreDropDown.AddSubMenuItem("Refresh", "")
	go RunRefresh(moreRefresh)

	mQuit := systray.AddMenuItem("Quit", "Quit this app")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
	}()
}

func CreateMenuRecursive(menus []FileConfig, menuParent *systray.MenuItem) {
	for _, item := range menus {
		if item.Divider {
			systray.AddSeparator()
		} else if item.Children != nil {
			option := systray.AddMenuItem(item.Label, "")
			CreateMenuRecursive(item.Children, option)
		} else {
			if menuParent != nil {
				option := menuParent.AddSubMenuItem(item.Label, "")
				go RunScript(option, item)
			} else {
				option := systray.AddMenuItem(item.Label, "")
				go RunScript(option, item)
			}
		}
	}
}

func CreateFileConfig() {
	if _, err := os.Stat(fileConfigPath); os.IsNotExist(err) {
		file, err := os.Create(fileConfigPath)
		if err != nil {
			println(&err)
			if err.Error() != "file exists" {
				log.Fatal(err)
			}
		}
		file.WriteString(yamlFileBase)
	} else {
		fmt.Println("The file exists:", fileConfigPath)
	}
}

func ReadFileConfig() RootConfig {
	var config RootConfig
	yamlFile, err := os.ReadFile(fileConfigPath)
	if err != nil {
		log.Fatalf(fileConfigPath+": err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return config
}

func RunScript(option *systray.MenuItem, item FileConfig) {
	for {
		<-option.ClickedCh
		fmt.Println("Run", item.Label)
		var cmd *exec.Cmd
		if runtime.GOOS == "darwin" {
			descOsAsScript := `tell application "Terminal" to do script "` + item.Run + `" activate`
			cmd = exec.Command("osascript", "-s", "h", "-e", descOsAsScript)
		} else if runtime.GOOS == "linux" {
			cmd = exec.Command("gnome-terminal", "--", "bash", "-c", item.Run)
		} else {
			log.Fatal("OS not supported")
		}
		stderr, err := cmd.StderrPipe()
		log.SetOutput(os.Stderr)
		if err != nil {
			println(err.Error())
		}
		if err := cmd.Start(); err != nil {
			println(err.Error())
		}
		slurp, _ := io.ReadAll(stderr)
		fmt.Printf("%s\n", slurp)
		if err := cmd.Wait(); err != nil {
			println(err.Error())
		}
		println("ok: ", item.Label)
	}
}

func RunRefresh(option *systray.MenuItem) {
	<-option.ClickedCh

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	filename := filepath.Base(os.Args[0])

	fmt.Println("Run refresh")
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		descOsAsScript := `tell application "Terminal" to do script "` + dir + "/" + filename + `" activate`
		cmd = exec.Command("osascript", "-s", "h", "-e", descOsAsScript)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("/usr/bin/nohup", dir+"/"+filename, "&")
	} else {
		log.Fatal("OS not supported")
	}
	if err := cmd.Start(); err != nil {
		println(err.Error())
	}

	println("ok: refresh")
	os.Exit(0)
}

var yamlFileBase = `
root:
- label: Hello World!
  run: echo Hello World!
`

var IconTerminal []byte = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x18,
	0x08, 0x06, 0x00, 0x00, 0x00, 0xe0, 0x77, 0x3d, 0xf8, 0x00, 0x00, 0x00,
	0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00, 0x0b, 0x13, 0x00, 0x00, 0x0b,
	0x13, 0x01, 0x00, 0x9a, 0x9c, 0x18, 0x00, 0x00, 0x00, 0xbe, 0x49, 0x44,
	0x41, 0x54, 0x78, 0x9c, 0xed, 0x95, 0x41, 0x0a, 0xc2, 0x30, 0x10, 0x45,
	0x73, 0x83, 0xf6, 0x1e, 0xda, 0xad, 0xdd, 0x08, 0x16, 0xef, 0xd2, 0x78,
	0x04, 0xb7, 0x5d, 0x7a, 0x8f, 0xd6, 0x53, 0xe9, 0x01, 0xac, 0x97, 0x78,
	0x12, 0x3a, 0xa0, 0xb6, 0xb1, 0xa4, 0x24, 0x03, 0x2e, 0xfa, 0x61, 0x36,
	0x13, 0xf2, 0x5f, 0xf8, 0x81, 0x19, 0x63, 0x56, 0x85, 0x0a, 0xc8, 0x01,
	0x0b, 0x9c, 0x23, 0xcb, 0x02, 0x99, 0xcf, 0xfc, 0x0e, 0xb4, 0xc0, 0x25,
	0xb2, 0x5a, 0xf1, 0x7a, 0x43, 0x80, 0x93, 0x3b, 0x48, 0x98, 0x46, 0x07,
	0xd4, 0x9f, 0x8d, 0xc6, 0x55, 0x42, 0x40, 0xf3, 0xe5, 0x37, 0x6e, 0x00,
	0x5b, 0x60, 0xa7, 0x09, 0x28, 0x81, 0x1e, 0x38, 0xa8, 0x00, 0xa4, 0xb7,
	0x17, 0xc8, 0x51, 0x05, 0x10, 0x03, 0x21, 0x14, 0x20, 0x67, 0x15, 0xf0,
	0x74, 0xb1, 0x69, 0x02, 0xfa, 0x25, 0x9f, 0xfe, 0x1f, 0x11, 0x31, 0x98,
	0x3f, 0x54, 0x3e, 0xd9, 0xc5, 0x21, 0x2f, 0xaf, 0x46, 0x17, 0x8b, 0x99,
	0x19, 0x54, 0x2c, 0x01, 0x6c, 0x7c, 0x99, 0x33, 0x00, 0x7e, 0xcd, 0xa0,
	0x70, 0x40, 0xac, 0xf0, 0x00, 0xac, 0xf6, 0xb0, 0xcb, 0x65, 0xc4, 0x76,
	0x09, 0xf6, 0xc1, 0x15, 0xb8, 0xf9, 0x76, 0x42, 0xe6, 0xa8, 0x09, 0xf6,
	0x41, 0x3d, 0x31, 0x5f, 0x65, 0x66, 0xf4, 0x02, 0xf0, 0x1a, 0x7c, 0x95,
	0x10, 0x6d, 0xf6, 0x7c, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
}
