package core

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/gflydev/core/utils"
	"io"
	"log"
)

func loadFile(template string) (*pongo2.Template, error) {
	viewPath := utils.Getenv("VIEW_PATH", "resources/views")
	viewExt := utils.Getenv("VIEW_EXT", "tpl")

	return pongo2.FromFile(fmt.Sprintf("%s/%s.%s", viewPath, template, viewExt))
}

// LoadView load view content from template.
func LoadView(template string, d Data) string {
	file, err := loadFile(template)
	if err != nil {
		log.Fatal(err)
	}
	ctx := pongo2.Context(d)
	res, err := file.Execute(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

// LoadViewWriter load view content into writer.
func LoadViewWriter(template string, d Data, w io.Writer) error {
	var templateFile = pongo2.Must(loadFile(template))
	ctx := pongo2.Context(d)

	return templateFile.ExecuteWriter(ctx, w)
}
