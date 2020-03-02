package main

import (
	"encoding/json"
	"flag"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	//import "github.com/Masterminds/squirrel"
	//import "github.com/hashicorp/go-getter"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func parse(dir, file string) *hclparse.Parser {
	parser := hclparse.NewParser()
	var files []string
	if dir != "" {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				if filepath.Ext(path) == ".hcl" {
					log.Debugf("Parser found hcl file: %s",path)
					files = append(files, path)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

	}
	if file != "" {
		info, err := os.Stat(file)
		if info.IsDir() {
			log.Fatalf("%s is a directory. Did you mean to use -dir ?", file)
		}
		if os.IsNotExist(err) {
			log.Fatalf("%s does not exist", file)
		}
	}
	for _, file := range files {
		f , diags := parser.ParseHCLFile(file)
		if diags != nil && diags.HasErrors() {
			log.Fatal(diags.Error())
		}
		log.Tracef("file:%s content:%s",file ,f.Bytes)
	}
	return parser
}

func main() {

	//Read command line arguments
	dirFlag := flag.String("dir", "", "parse .hcl files in directory")
	fileFlag := flag.String("file", "", "parse file")
	flag.Parse()
	log.Debug("CLI args:")
	if *dirFlag != "" {
		log.Debugf("-dir=%s",*dirFlag)
	}
	if *fileFlag != "" {
		log.Debugf("-file %s",*fileFlag)
	}

	// Parse hcl files into single body
	parser := parse(*dirFlag, *fileFlag)
	var files []*hcl.File
	for _, file := range parser.Files() {
		files = append(files, file)
	}
	body := hcl.MergeFiles(files)

	// Decode into Dvml struct
	var dvml Dvml
	diag := gohcl.DecodeBody(body,nil, &dvml )
	if diag != nil {
		log.Fatal(diag.Error())
	}


	s, _ := json.MarshalIndent(dvml,"","\t")
	log.Debug(string(s))

}

type Dvml struct {
	Source Source `hcl:"source,block"`
	Target Target `hcl:"target,block"`
	//Remain hcl.Body `hcl:"denna,remain"`
}

type Source struct {
	Json Json `hcl:"json,block"`
}

type Json struct {
	Name string `hcl:"name,label"`
	Fields Fields `hcl:"fields,block" json:"fields"`
}

type Fields struct {
	Varchar []Varchar `hcl:"varchar,block"`
}

type Varchar struct {
	Name string `hcl:"name,label"`
	Path string `hcl:"path,attr" `
}

type Target struct {
	Hub []Hub `hcl:"hub,block"`
}
type Hub struct {
	Name string `hcl:"name,label"`
	Key string `hcl:"key,attr"`
}
