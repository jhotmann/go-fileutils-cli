package operation

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/flosch/pongo2/v4"
	_ "github.com/jhotmann/go-fileutils-cli/lib/filters"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

type Operation struct {
	Type           string
	Input          util.PathObject
	Stats          os.FileInfo
	OutputTemplate *pongo2.Template
	Output         util.PathObject
}

type OperationList []Operation

func FilesToOperationsList(opType string, files []string, outputTemplate *pongo2.Template) OperationList {
	operations := []Operation{}
	for _, f := range files {
		matches, err := filepath.Glob(f)
		if err != nil {
			panic(err)
		}
		for _, match := range matches {
			var op Operation
			op.Input = util.GetPathObj(match)
			op.OutputTemplate = outputTemplate
			stats, err := os.Stat(match)
			if err != nil {
				fmt.Println(err)
				continue
			}
			op.Stats = stats
			operations = append(operations, op)
		}
	}
	return operations
}

func (o OperationList) RemoveDirectories() OperationList {
	directoryIndex := []int{}
	for i, op := range o {
		if op.Stats.IsDir() {
			directoryIndex = append(directoryIndex, i)
		}
	}
	for i, dirIndex := range directoryIndex { // I don't know why this works, but it does
		if dirIndex == len(o)-1 {
			return o[:dirIndex-i]
		} else {
			return append(o[:dirIndex-i], o[dirIndex+1-i])
		}
	}
	return o
}

func (o OperationList) RemoveDuplicateInputs() OperationList {
	ret := OperationList{}
	paths := []string{}
	for _, op := range o {
		abs := op.Input.Abs
		if util.IndexOf(abs, paths) == -1 {
			paths = append(paths, abs)
			ret = append(ret, op)
		}
	}
	return ret
}

func (o OperationList) Sort(sortOption string) OperationList {
	switch sortOption {
	case "alphabet":
		sort.Slice(o, func(i, j int) bool { return o[i].Input.Abs < o[j].Input.Abs })
	case "reverse-alphabet":
		sort.Slice(o, func(i, j int) bool { return o[j].Input.Abs < o[i].Input.Abs })
	case "date":
		sort.Slice(o, func(i, j int) bool { return o[i].Stats.ModTime().Before(o[j].Stats.ModTime()) })
	case "reverse-date":
		sort.Slice(o, func(i, j int) bool { return o[i].Stats.ModTime().After(o[j].Stats.ModTime()) })
	case "size":
		sort.Slice(o, func(i, j int) bool { return o[i].Stats.Size() > o[j].Stats.Size() })
	case "reverse-size":
		sort.Slice(o, func(i, j int) bool { return o[j].Stats.Size() > o[i].Stats.Size() })
	}
	return o
}

func (o OperationList) RenderTemplates() OperationList {
	ret := OperationList{}
	for _, op := range o {
		context := pongo2.Context{
			"i":           "--FILEINDEXHERE--",
			"f":           op.Input.Name,
			"abs":         op.Input.Abs,
			"rel":         op.Input.Rel,
			"ext":         op.Input.Ext,
			"p":           filepath.Dir(op.Input.Dir),
			"isDirectory": fmt.Sprintf("%t", op.Stats.IsDir()),
			"date": map[string]time.Time{
				"now":      time.Now(),
				"modified": op.Stats.ModTime(),
			},
			"size": op.Stats.Size(),
		}
		out, err := op.OutputTemplate.Execute(context)
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
		op.Output = util.GetPathObj(out)
		ret = append(ret, op)
	}
	return ret
}
