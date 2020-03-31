package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func Parse(dir string) *hclparse.Parser {
	parser := hclparse.NewParser()
	var files []string
	if dir != "" {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				if filepath.Ext(path) == ".hcl" {
					log.Debugf("Parser found hcl file: %s", path)
					files = append(files, path)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, file := range files {
		f, diags := parser.ParseHCLFile(file)
		if diags != nil && diags.HasErrors() {
			log.Fatal(diags.Error())
		}
		log.Tracef("file:%s content:%s", file, f.Bytes)
	}
	return parser
}

type Source interface {
	Debug()
}

// PetsHCL is a generic structure that could be either cats or dogs. The Type
// field indicates which, and the generic "characteristics" block HCL will be
// decoded into the unique fields for that type.
// Note the use of the `hcl:",remain"` tag, which puts all undecoded HCL into
// an hcl.Body for use later.
type SourcesHCL struct {
	SourceHCLBodies []*struct {
		Name          string `hcl:",label"`
		Type          string `hcl:"type"`
		AttributesHCL *struct {
			HCL hcl.Body `hcl:",remain"`
		} `hcl:"attributes,block"`
	} `hcl:"source,block"`
}

type Json struct {
	Varchar string `hcl:"varchar,optional"`
}

func (s *Json) Debug() {
	fmt.Printf("varchar: %s", s.Varchar)
}

func ParseHCL(dir string) ([]Source, error) {
	parser := Parse(dir)
	var files []*hcl.File
	for _, file := range parser.Files() {
		files = append(files, file)
	}

	body := hcl.MergeFiles(files)

	evalContext := SourceEvalContext()

	sourcesHCL := &SourcesHCL{}

	if diag := gohcl.DecodeBody(body, evalContext, sourcesHCL); diag.HasErrors() {
		return []Source{}, fmt.Errorf("error in decode HCL: %w", diag)
	}

	sources := []Source{}

	for _, s := range sourcesHCL.SourceHCLBodies {
		switch sourceType := s.Type; sourceType {
		case "json":
			json := &Json{}
			if s.AttributesHCL != nil {
				if diag := gohcl.DecodeBody(s.AttributesHCL.HCL, evalContext, json); diag.HasErrors() {
					return []Source{}, fmt.Errorf("error in decode HCL: %w", diag)
				}
			}
			sources = append(sources, json)
		default:
			return []Source{}, fmt.Errorf("error in decoding: unknown source type `%s`", sourceType)
		}
	}
	return sources, nil
}

func SourceEvalContext() *hcl.EvalContext {
	return nil
}
