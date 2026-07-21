package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type GameInfo struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("pack it up bro pls use me correctly")
	}

	root := os.Args[1]

	scripts := filepath.Join(root, "scripts")
	styles := filepath.Join(root, "styles")
	games := filepath.Join(root, "games")

	out := filepath.Join(root, "dist")

	build(api.BuildOptions{
		EntryPoints: []string{filepath.Join(scripts, "nes.js")},
		External:    []string{"*.png"},
		Bundle:      true,
		Splitting:   true,
		Format:      api.FormatESModule,
		Outdir:      out,
		Write:       true,
	})

	for _, name := range []string{"emuWorker.js", "driverWorklet.js"} {
		build(api.BuildOptions{
			EntryPoints: []string{filepath.Join(scripts, name)},
			Bundle:      true,
			Outdir:      out,
			Write:       true,
		})
	}

	build(api.BuildOptions{
		EntryPoints:       []string{filepath.Join(styles, "all.css")},
		External:          []string{"*.png"},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Outdir:            out,
		Write:             true,
	})

	fmt.Println("build finished")

	gamesList := make([]GameInfo, 0)
	filepath.WalkDir(games, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := filepath.Base(path)
		id, ok := strings.CutSuffix(name, ".nes")

		if !ok {
			return nil
		}

		gamesList = append(gamesList, GameInfo{
			ID:   id,
			Name: strings.ToTitle(strings.ReplaceAll(id, "_", " ")),
		})

		return nil
	})

	out = filepath.Join(root, "games.json")
	f, err := os.Create(out)

	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()
	if err := json.NewEncoder(f).Encode(gamesList); err != nil {
		log.Fatalln(err)
	}

}

func build(opts api.BuildOptions) {
	res := api.Build(opts)
	for _, err := range res.Errors {
		log.Fatalln(err)
	}
}
