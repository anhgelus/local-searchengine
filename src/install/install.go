package install

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/fatih/color"
)

type appInfo struct {
	Path string
}

type Configuration struct {
	Version       string
	AppName       string
	BlockList     []string
	WallpaperPath string
}

//go:embed world.anhgelus.local-searchengine.plist
var macOSService string

//go:embed local-searchengine.service
var linuxService string

func App() error {
	var err error
	nixUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("impossible de récupérer l'utilisateur courant %w", err)
	}
	home := nixUser.HomeDir
	config, err := toml.Marshal(Configuration{
		AppName:       "Local SearchEngine",
		Version:       "0.1",
		BlockList:     []string{""},
		WallpaperPath: "",
	})
	if err != nil {
		return fmt.Errorf("impossible de générer la configuration %s", err)
	}
	switch runtime.GOOS {
	case "darwin":
		plistPath := filepath.Join(home, "Library/LaunchAgents/world.anhgelus.local-searchengine.plist")
		err = installService(
			macOSService,
			plistPath,
			exec.Command("launchctl", "load", plistPath),
		)
		// TODO: generate the configuration for MacOS
		if err == nil {
			color.Green("Le service a été installé dans %s et activé !\n", plistPath)
			fmt.Println("")
			fmt.Println("Pour le désactiver :")
			color.Blue("launchctl unload %s", plistPath)
		}
	case "linux":
		linuxPath := filepath.Join(home, ".config/systemd/user/local-searchengine.service")
		err = installService(
			linuxService,
			linuxPath,
			exec.Command("systemctl", "enable", "--user", "local-searchengine.service"),
		)
		configPath := filepath.Join(home, ".config/local-searchengine/config.toml")
		err = os.WriteFile(configPath, config, 0644)
		if err != nil {
			return fmt.Errorf("impossible de créer et/ou écrire sur le fichier de configuration %s", err)
		}
		color.Green("Le service a été installé dans %s et activé !\n", linuxPath)
		fmt.Println("")
		fmt.Println("Pour le démarrer :")
		color.Blue("systemctl start --user local-searchengine.service")
		fmt.Println("")
		fmt.Println("Pour le désactiver :")
		color.Blue("systemctl disable --nixUser local-searchengine.service")
	default:
		return fmt.Errorf("système d'exploitation non géré %s", runtime.GOOS)
	}
	return err
}

func installService(tpl string, dest string, cmd *exec.Cmd) error {
	exePath, err := getCurrentExecPath()
	if err != nil {
		return err
	}

	t, err := template.New("service").Parse(tpl)
	if err != nil {
		return fmt.Errorf("impossible de parser le template %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier de service %w", err)
	}

	err = t.Execute(f, appInfo{Path: exePath})
	if err != nil {
		return err
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("impossible de démarrer le service %s %w", stderr.String(), err)
	}
	return nil
}

func getCurrentExecPath() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
