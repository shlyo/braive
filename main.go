//go:generate go install -v github.com/kevinburke/go-bindata/go-bindata
//go:generate go-bindata -prefix res/ -pkg assets -o assets/assets.go res/Brave.lnk
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"path"

	"github.com/shlyo/kernal"
	"github.com/shlyo/kernal/pkg/utl"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: true,
	}

	// Init app
	if app, err = portapps.NewWithCfg("brave", "Brave", cfg); err != nil {
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "brave.exe")
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
		"--disable-brave-update",
		"--no-default-browser-check",
		"--disable-logging",
		"--disable-breakpad",
		"--disable-machine-id",
		"--disable-encryption-win",
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			utl.Cleanup([]string{
				path.Join(os.Getenv("APPDATA"), "BraveSoftware"),
				path.Join(os.Getenv("LOCALAPPDATA"), "BraveSoftware"),
			})
		}()
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
