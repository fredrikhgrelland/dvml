package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
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
	log.SetLevel(log.InfoLevel)
}

const (
	dir = "conf"
)

func inner() error {
	//Read command line arguments
	dirFlag := flag.String("dir", dir, "parse .hcl files in directory")
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("CLI args:")

	if *dirFlag != "" {
		log.Debugf("-dir=%s", *dirFlag)
	}

	sources, err := ParseHCL(*dirFlag)
	if err != nil {
		return err
	}

	for _, s := range sources {
		s.Debug()
	}
	return nil
}

func main() {
	if err := inner(); err != nil {
		fmt.Printf("dvml error: %s\n", err.Error())
		os.Exit(1)
	}
}

/*
// Parse hcl files into single body
parser := parse(*dirFlag, *fileFlag)
var files []*hcl.File
for _, file := range parser.Files() {
files = append(files, file)
}
body := hcl.MergeFiles(files)

// Decode into Root struct
// Hierarki source -> target(hub) -> target(link) -> target(sat)
var root Root
diag := gohcl.DecodeBody(body, nil, &root)
if diag != nil {
for _, diag := range diag {
fmt.Printf("decoding - %s\n", diag)
}
log.Fatal(diag.Error())
}

variables := map[string]cty.Value{}

for _, v := range root.Source[0].Json.Fields.Varchar {
if len(v.Path) == 0 {
continue
}

val, diag := v.Path["path"].Expr.Value(nil)

if diag != nil {
for _, diag := range diag {
fmt.Printf("decoding - %s\n", diag)
}
log.Fatal(diag.Error())
}
variables[v.Name] = val
}

evalContext := &hcl.EvalContext{
Variables: map[string]cty.Value{
"var": cty.ObjectVal(variables),
"nows": cty.StringVal(time.Now().Format(time.RFC3339)),
},
Functions: map[string]function.Function{
"upper": stdlib.UpperFunc,
"now": now(),
},
}

// create output definition
spec := &hcldec.ObjectSpec{
"key": &hcldec.AttrSpec{
Name:     "key",
Required: true,
Type:     cty.String,
},
"date": &hcldec.AttrSpec{
Name:     "date",
Required: false,
Type:     cty.String,

},
"computed": &hcldec.AttrSpec{
Name:     "computed",
Required: false,
Type:     cty.String,

},
}

for _, target := range root.Target {
for _, hub := range target.Hub {
cfg, diag := hcldec.Decode(hub.Config, spec, evalContext)
if diag != nil {
for _, diag := range diag {
fmt.Printf("decoding - %s\n", diag)
}
log.Fatal(diag.Error())
}
result, _ := json.MarshalIndent(ctyjson.SimpleJSONValue{cfg}, "", " ")
log.Info(string(result))
}
}

s, _ := json.MarshalIndent(root, "", " ")
log.Debug(string(s))

type Root struct {
	Source []Source `hcl:"source,block"`
	Target []Target `hcl:"target,block"`
}

type Source struct {
	Json Json `hcl:"json,block"`
}

type Json struct {
	Name   string `hcl:"name,label"  json:"name"`
	Fields Fields `hcl:"fields,block" json:"fields"`
}

type Fields struct {
	Varchar []Varchar `hcl:"varchar,block"`
	Numbers []Number  `hcl:"number,block"`
}

type Varchar struct {
	Name string         `hcl:"name,label"`
	Path hcl.Attributes `hcl:"path,remain"`
}
type Number struct {
	Name string `hcl:"name,label"`
	Path string `hcl:"path,attr" `
}

type Target struct {
	Hub    []Hub    `hcl:"hub,block"`
}

type Hub struct {
	Name string `hcl:"name,label"`
	Computed string `hcl:"computed,optional"`
	Config hcl.Body `hcl:",remain"`
}

// Now!
func now() function.Function {
	return function.New(&function.Spec{
		Params: nil,
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			return cty.StringVal(time.Now().Format(time.RFC3339)), nil
		},
	})
}
*/
