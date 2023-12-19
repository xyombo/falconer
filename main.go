package main

import (
	"fmt"
	"github.com/rivo/tview"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Desc     string `yaml:"desc"`
	Password string `yaml:"password"`
}

type YamlConfig struct {
	Servers []ServerConfig `yaml:"servers"`
}

func LoadServerConfig(file string) ([]ServerConfig, error) {
	var config YamlConfig
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	// 解析 YAML 内容
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}
	// 输出解析后的配置
	return config.Servers, nil
}

type SelectRowState struct {
	index int
}

func ShowSelectMenu(serverConfigs []ServerConfig) (ServerConfig, error) {
	app := tview.NewApplication()
	list := tview.NewList()
	state := SelectRowState{}

	for i, s := range serverConfigs {
		list.AddItem(s.Host, s.Desc, rune(i), nil)
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		state.index = index
		app.Stop()
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
	return serverConfigs[state.index], nil
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// 构建文件路径（这里以读取 .bashrc 文件为例）
	filePath := filepath.Join(homeDir, ".ssh/stabby_config.yaml")
	// 加载配置
	serverConfigs, _ := LoadServerConfig(filePath)

	//把配置的服务器列表用tview.list展示出来
	if serverConfig, err := ShowSelectMenu(serverConfigs); err == nil {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), &ssh.ClientConfig{
			User:            "root",
			Auth:            []ssh.AuthMethod{ssh.Password(serverConfig.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})

		session, err := client.NewSession()
		defer func(session *ssh.Session) {
			err := session.Close()
			if err != nil {

			}
		}(session)
		if err != nil {
			log.Fatalf("new session error: %s", err.Error())
		}

		// 获取伪终端
		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			log.Fatalf("terminal make raw: %s", err)
		}
		defer func(fd int, oldState *terminal.State) {
			err := terminal.Restore(fd, oldState)
			if err != nil {

			}
		}(fd, state)

		session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
		session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
		session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,     // 禁用回显（0禁用，1启动）
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
		}
		if err = session.RequestPty("xterm", 32, 160, modes); err != nil {
			log.Fatalf("request pty error: %s", err.Error())
		}
		if err = session.Shell(); err != nil {
			log.Fatalf("start shell error: %s", err.Error())
		}
		if err = session.Wait(); err != nil {
			log.Fatalf("return error: %s", err.Error())
		}
		if err != nil {
			log.Fatalf("SSH dial error: %s", err.Error())
		}
	}

}
