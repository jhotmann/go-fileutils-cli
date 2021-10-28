package operation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1set/gut/yos"
	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/db"
	"github.com/jhotmann/go-fileutils-cli/util"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

func init() {
	pongo2.ReplaceFilter("date", DateFilter)
	pongo2.ReplaceFilter("time", DateFilter)
	pongo2.ReplaceFilter("title", titleFilter)
	pongo2.RegisterFilter("pascal", titleFilter)
	pongo2.RegisterFilter("snake", snakeFilter)
	pongo2.RegisterFilter("camel", camelFilter)
	pongo2.RegisterFilter("kebab", kebabFilter)
	pongo2.RegisterFilter("replace", replaceFilter)
	pongo2.RegisterFilter("regexReplace", regexReplaceFilter)
	pongo2.RegisterFilter("with", withFilter)
	pongo2.RegisterFilter("match", matchFilter)
	pongo2.RegisterFilter("index", indexFilter)
	pongo2.RegisterFilter("pad", padFilter)
}

type Operation struct {
	Type           string
	Input          PathObject
	Stats          os.FileInfo
	OutputTemplate *pongo2.Template
	Output         PathObject
	HasConflict    bool
	Index          int
	ConflictCount  int
	Options        OperationOptions
	Skip           bool
}

type OperationList []Operation

var OperationType = struct {
	Mv string
	Cp string
	Ln string
}{
	Mv: "move",
	Cp: "copy",
	Ln: "link",
}

func (opl OperationList) Initialize() OperationList {
	inputs := []string{}
	counts := map[string]int{}
	for i, op := range opl {
		if op.Options.IgnoreDirectories && op.Stats.IsDir() { // Ignore directories if specified
			op.Skip = true
		}
		if util.IndexOf(op.Input.Abs, inputs) == -1 { // If input is unique
			inputs = append(inputs, op.Input.Abs)
		} else { // Ignore duplicate inputs
			op.Skip = true
		}
		if op.Skip {
			opl[i] = op
			continue
		}
		op.Output = op.renderTemplate()
		if !op.Options.NoExt && op.Input.Ext != "" && op.Output.Ext == "" { // Add extension if not specified and not --no-ext used
			op.Output = op.Output.UpdateExt(op.Input.Ext)
		}
		if op.Options.NoMove { // Don't move files if --no-move option used
			op.Output = op.Output.UpdateDir(op.Input.Dir)
		}
		count := counts[op.Output.Abs]
		counts[op.Output.Abs] = count + 1
		op.Index = count + 1
		if count == 0 || op.Options.Force {
			op.HasConflict = false
		} else {
			op.HasConflict = true
		}
		opl[i] = op
	}

	for i, op := range opl { // Loop through a second time to set indexes
		if !op.Options.NoIndex && (op.HasConflict || counts[op.Output.Abs] > 1) {
			index := util.ZeroPad(op.Index, counts[op.Output.Abs])
			if strings.Contains(op.Output.Name, "--FILEINDEXHERE--") { // put index where {{i}} was used
				op.Output = op.Output.UpdateName(strings.Replace(op.Output.Name, "--FILEINDEXHERE--", index, -1))
			} else { // put index at end of file name
				op.Output = op.Output.UpdateName(op.Output.Name + index)
			}
		} else { // remove any --FILEINDEXHERE-- if {{i}} was accidentally used
			op.Output = op.Output.UpdateName(strings.Replace(op.Output.Name, "--FILEINDEXHERE--", "", -1))
		}
		opl[i] = op
	}
	return opl
}

func (opl OperationList) Run(command []string) {
	var batch db.Batch
	defer batch.Close()

	if len(opl) > 0 && !opl[0].Options.Simulate {
		batch = db.NewBatch(opl[0].Type, command, util.GetWorkingDir())
	}

	for i, op := range opl {
		var err error
		if op.Skip {
			continue
		}
		if op.Options.Simulate {
			pterm.Info.Printfln("%s → %s", op.Input.Rel, op.Output.Rel)
			continue
		}
		if op.Input.Abs == op.Output.Abs { // no change
			if op.Options.Verbose {
				pterm.Info.Printfln("Skipping %s because it did not change", op.Input.Rel)
			}
			continue
		}
		if !op.Options.NoMkdir { // Make sure all output directories exist
			_, err := os.Stat(op.Output.Dir)
			if os.IsNotExist(err) { // Output directory doesn't exist so we'll create it with the same permissions as the input file
				stats, _ := os.Stat(op.Input.Dir)
				os.MkdirAll(op.Output.Dir, stats.Mode()) // Make output directory with same permissions as input directory
			}
		}
		if op.Options.Force { // Do the operation with reckless abadon
			err = op.runOperation()
		} else {
			_, existErr := os.Stat(op.Output.Abs)
			if os.IsNotExist(existErr) { // File/Dir doesn't exist so we can proceed
				err = op.runOperation()
			} else { // File/Dir already exists, check with user what to do
				if strings.EqualFold(op.Input.Abs, op.Output.Abs) && op.Type == OperationType.Mv { // rename with case change, allow it
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
					case 1:
						prompt2 := promptui.Prompt{
							Label:   "New File Name",
							Default: op.Output.Name,
						}
						val, err := prompt2.Run()
						if err != nil {
							panic(err)
						}
						opl[i].Output = opl[i].Output.UpdateName(val)
						opl[i].runOperation()
					case 2:
						if op.Options.Verbose {
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
		if op.Options.Verbose {
			pterm.Success.Printfln("%s → %s", op.Input.Rel, op.Output.Rel)
		}
	}
}

func (op Operation) renderTemplate() PathObject {
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
	return GetPathObj(out)
}

func (o Operation) runOperation() error {
	var err error
	switch o.Type {
	case OperationType.Mv:
		err = os.Rename(o.Input.Abs, o.Output.Abs)
	case OperationType.Cp:
		if o.Stats.IsDir() {
			yos.CopyDir(o.Input.Abs, o.Output.Abs)
		} else {
			yos.CopyFile(o.Input.Abs, o.Output.Abs)
		}
	case OperationType.Ln:
		if o.Options.Soft {
			err = os.Symlink(o.Input.Abs, o.Output.Abs)
		} else {
			err = os.Link(o.Input.Abs, o.Output.Abs)
		}
	default:
		err = errors.New(o.Type + " not implemented")
	}
	return err
}
