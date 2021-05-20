package operation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/1set/gut/yos"
	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/lib/db"
	_ "github.com/jhotmann/go-fileutils-cli/lib/filters"
	"github.com/jhotmann/go-fileutils-cli/lib/options"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

type Operation struct {
	Type           string
	Input          util.PathObject
	Stats          os.FileInfo
	OutputTemplate *pongo2.Template
	Output         util.PathObject
	HasConflict    bool
	Index          int
	ConflictCount  int
}

type OperationList []Operation

func FilesToOperationsList(opType string, files []string, outputTemplate *pongo2.Template) OperationList {
	operations := []Operation{}
	for _, f := range files {
		matches, err := filepath.Glob(f)
		if err != nil {
			panic(err)
		}
		if len(matches) == 0 {
			pterm.Warning.Printfln("%s does not match any existing files", f)
		}
		for _, match := range matches {
			var op Operation
			op.Type = opType
			op.Input = util.GetPathObj(match)
			op.OutputTemplate = outputTemplate
			stats, err := os.Stat(match)
			if err != nil {
				pterm.Warning.Println(err.Error())
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
		out = strings.ReplaceAll(out, "--REPLACEME--", "")
		op.Output = util.GetPathObj(out)
		ret = append(ret, op)
	}
	return ret
}

func (o OperationList) PopulateBlankExtensions() OperationList {
	ret := OperationList{}
	for _, op := range o {
		if op.Output.Ext == "" && op.Input.Ext != "" {
			op.Output = op.Output.UpdateExt(op.Input.Ext)
		}
		ret = append(ret, op)
	}
	return ret
}

func (o OperationList) NoMove() OperationList {
	ret := OperationList{}
	for _, op := range o {
		op.Output = op.Output.UpdateDir(op.Input.Dir)
		ret = append(ret, op)
	}
	return ret
}

func (o OperationList) FindConflicts() OperationList {
	ret := OperationList{}
	counts := map[string]int{}
	indicies := map[string]int{}
	for _, op := range o {
		value := counts[op.Output.Abs] // If key doesn't exist, value will be zero
		counts[op.Output.Abs] = value + 1
	}
	for _, op := range o {
		count := counts[op.Output.Abs]
		op.ConflictCount = count
		if count == 1 {
			op.Index = 1
			op.HasConflict = false
		} else {
			op.HasConflict = true
			index := indicies[op.Output.Abs] + 1
			indicies[op.Output.Abs] = index
			op.Index = index
		}
		ret = append(ret, op)
	}
	return ret
}

func (o OperationList) AddIndex() OperationList {
	ret := OperationList{}
	for _, op := range o {
		if op.HasConflict {
			index := util.ZeroPad(op.Index, op.ConflictCount)
			if strings.Contains(op.Output.Name, "--FILEINDEXHERE--") {
				op.Output = op.Output.UpdateName(strings.Replace(op.Output.Name, "--FILEINDEXHERE--", index, -1))
			} else {
				op.Output = op.Output.UpdateName(op.Output.Name + index)
			}
		} else {
			op.Output = op.Output.UpdateName(strings.Replace(op.Output.Name, "--FILEINDEXHERE--", "", -1))
		}
		ret = append(ret, op)
	}
	return ret
}

func (o OperationList) Run(command []string, opts options.CommonOptions) {
	var batch db.Batch
	if !opts.Simulate {
		batch = db.NewBatch(o[0].Type, command, util.GetWorkingDir())
	}
	defer batch.Close()
	for _, op := range o {
		var err error
		if opts.Simulate {
			pterm.Info.Printfln("%s → %s", op.Input.Rel, op.Output.Rel)
			continue
		}
		if op.Input.Abs == op.Output.Abs { // no change
			if opts.Verbose {
				pterm.Info.Printfln("Skipping %s because it did not change", op.Input.Rel)
			}
			continue
		}
		if !opts.NoMkdir { // Make sure all output directories exist
			_, err := os.Stat(op.Output.Dir)
			if os.IsNotExist(err) { // Output directory doesn't exist so we'll create it with the same permissions as the input file
				stats, _ := os.Stat(op.Input.Dir)
				os.MkdirAll(op.Output.Dir, stats.Mode())
			}
		}
		if opts.Force { // Do the operation with reckless abadon
			err = op.runOperation()
		} else {
			_, err := os.Stat(op.Output.Abs)
			if os.IsNotExist(err) { // File/Dir doesn't exist so we can proceed
				err = op.runOperation()
			} else { // File/Dir already exists, check with user what to do
				if strings.ToLower(op.Input.Abs) == strings.ToLower(op.Output.Abs) && op.Type == "move" { // rename with case change, allow it
					op.runOperation()
				} else { // Prompt for user input
					fmt.Println()
					pterm.Warning.Printfln("What should happen to %s, %s already exists", op.Input.Rel, op.Output.Rel)
					prompt := promptui.Select{
						Label: "What would you like to do?",
						Items: []string{"Overwrite", "Input a new name", "Skip"},
					}
					index, _, _ := prompt.Run()
					switch index {
					case 0:
						op.runOperation()
						break
					case 1:
						prompt2 := promptui.Prompt{
							Label:   "New File Name",
							Default: op.Output.Name,
						}
						val, err := prompt2.Run()
						if err != nil {
							panic(err)
						}
						op.Output = op.Output.UpdateName(val)
						op.runOperation()
						break
					case 2:
						if opts.Verbose {
							pterm.Info.Printfln("Skipping %s", op.Input.Rel)
						}
						continue
					}
				}
			}
		}
		if err != nil {
			pterm.Error.Println(err.Error())
			continue
		}
		db.WriteOperation(batch.Id, op.Input.Abs, op.Output.Abs)
		if opts.Verbose {
			pterm.Success.Printfln("%s → %s", op.Input.Rel, op.Output.Rel)
		}
	}
}

func (o Operation) runOperation() error {
	var err error
	switch o.Type {
	case "move":
		err = os.Rename(o.Input.Abs, o.Output.Abs)
	case "copy":
		if o.Stats.IsDir() {
			yos.CopyDir(o.Input.Abs, o.Output.Abs)
		} else {
			yos.CopyFile(o.Input.Abs, o.Output.Abs)
		}
	case "link-soft":
		err = os.Symlink(o.Input.Abs, o.Output.Abs)
	case "link-hard":
		err = os.Link(o.Input.Abs, o.Output.Abs)
	default:
		err = errors.New(o.Type + " not implemented")
	}
	return err
}
