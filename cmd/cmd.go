// Package cmd implements the glx-web command.
package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/internal/load"
	"github.com/dhemery/glx-web/internal/render"
	"github.com/spf13/cobra"
)

func Execute() {
	err := webCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var webCmd = &cobra.Command{
	Use:   "glx-web",
	Short: "Build a website from a GLX archive and templates.",
	Args:  cobra.NoArgs,
	RunE:  run,
}

var (
	archiveDir          = "."
	buildSite           = true
	cleanBeforeBuilding = false
	outputDir           = "./public"
	serveAddress        = ""
	staticDir           = ""
	templateDir         = "./templates"
)

func init() {
	webCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"the `archive` to render")
	webCmd.Flags().BoolVarP(&buildSite, "build", "b", buildSite,
		"build the website")
	webCmd.Flags().BoolVarP(&cleanBeforeBuilding, "clean", "c", cleanBeforeBuilding,
		"remove the existing output dir before building")
	webCmd.Flags().StringVarP(&outputDir, "output", "o", outputDir,
		"the output `dir`")
	webCmd.Flags().StringVarP(&serveAddress, "serve", "S", serveAddress,
		"serve the website on `address`, e.g. \"localhost:3333\"")
	webCmd.Flags().StringVarP(&staticDir, "static", "s", staticDir,
		"a `dir` of files to copy into the output dir, if any")
	webCmd.Flags().StringVarP(&templateDir, "templates", "t", templateDir,
		"the template `dir`")

}

func run(_ *cobra.Command, _ []string) error {
	if buildSite {
		if err := build(); err != nil {
			return err
		}
	}
	if serveAddress != "" {
		return serve()
	}
	return nil
}

func build() error {
	absArchiveDir, err := filepath.Abs(archiveDir)
	if err != nil {
		return fmt.Errorf("archive directory: %w", err)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("output directory: %w", err)
	}

	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return fmt.Errorf("template directory: %w", err)
	}

	glxFile, err := load.GLX(absArchiveDir)
	if err != nil {
		return err
	}

	templates, err := load.Templates(absTemplateDir)
	if err != nil {
		return err
	}

	if cleanBeforeBuilding {
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("removing output dir: %w", err)
		}
	}

	if err := os.Mkdir(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if staticDir != "" {
		if err := os.CopyFS(outputDir, os.DirFS(staticDir)); err != nil {
			return fmt.Errorf("copying static files: %w", err)
		}
	}

	r := &render.Renderer{
		OutputDir: absOutputDir,
		Templates: templates,
	}

	catalog := entity.NewCatalog(glxFile, os.Stderr)

	if err := r.Render(catalog, glxFile); err != nil {
		return fmt.Errorf("rendering: %w", err)
	}
	return nil
}

func serve() error {
	listener, err := net.Listen("tcp", serveAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	http.Handle("/", http.FileServer(http.Dir(outputDir)))

	fmt.Println("Server is running on http://" + listener.Addr().String())

	return http.Serve(listener, nil)
}
