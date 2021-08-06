//+build generate

// gen_entity.go generates entity information.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
)

const (
	infoURL = "https://raw.githubusercontent.com/PrismarineJS/minecraft-data/master/data/pc/1.17/entities.json"
	//language=gohtml
	entityTmpl = `// Code generated by gen_entity.go DO NOT EDIT.
// Package entity stores information about entities in Minecraft.
package entity
// ID describes the numeric ID of an entity.
type ID uint32

// Entity describes information about a type of entity.
type Entity struct {
ID          ID
InternalID  uint32
DisplayName string
Name        string
Width  float64
Height float64
Type     string
}

var (
	{{- range .}}
	{{.CamelName}} = Entity{
		ID: {{.ID}},
		InternalID: {{.InternalID}},
		DisplayName: "{{.DisplayName}}",
		Name: "{{.Name}}",
		Width: {{.Width}},
		Height: {{.Height}},
		Type: "{{.Type}}",
	}{{end}}
)

// ByID is an index of minecraft entities by their ID.
var ByID = map[ID]*Entity{ {{range .}}
	{{.ID}}: &{{.CamelName}},{{end}}
}`
)

type Entity struct {
	ID          uint32 `json:"id"`
	InternalID  uint32 `json:"internalId"`
	CamelName   string `json:"-"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`

	Width  float64 `json:"width"`
	Height float64 `json:"height"`

	Type string `json:"type"`
}

func downloadInfo() ([]*Entity, error) {
	resp, err := http.Get(infoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []*Entity
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	for _, d := range data {
		d.CamelName = strcase.ToCamel(d.Name)
	}
	return data, nil
}

//go:generate go run $GOFILE
//go:generate go fmt entity.go
func main() {
	fmt.Println("generating entity.go")
	entities, err := downloadInfo()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("entity.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := template.Must(template.New("").Parse(entityTmpl)).Execute(f, entities); err != nil {
		panic(err)
	}
}
