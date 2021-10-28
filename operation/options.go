package operation

import (
	"sort"
)

type OperationOptions struct {
	Force             bool
	Simulate          bool
	Sort              string
	Verbose           bool
	IgnoreDirectories bool
	NoIndex           bool
	NoExt             bool
	NoMkdir           bool
	NoMove            bool
	Soft              bool
}

var DefaultOptions = OperationOptions{
	Force:             false,
	Simulate:          false,
	Sort:              SortType.Alphabet,
	Verbose:           false,
	IgnoreDirectories: false,
	NoIndex:           false,
	NoExt:             false,
	NoMkdir:           false,
	NoMove:            false,
	Soft:              false,
}

var SortType = struct {
	None            string
	Alphabet        string
	ReverseAlphabet string
	Date            string
	ReverseDate     string
	Size            string
	ReverseSize     string
}{
	None:            "none",
	Alphabet:        "alphabet",
	ReverseAlphabet: "reverse-alphabet",
	Date:            "date",
	ReverseDate:     "reverse-date",
	Size:            "size",
	ReverseSize:     "reverse-size",
}

var AllowedSortValues = []string{
	SortType.None,
	SortType.Alphabet,
	SortType.ReverseAlphabet,
	SortType.Date,
	SortType.ReverseDate,
	SortType.Size,
	SortType.ReverseSize,
}

func (opl OperationList) WithForce(force bool) OperationList {
	for i := range opl {
		opl[i].Options.Force = force
	}
	return opl
}

func (opl OperationList) WithSimulate(simulate bool) OperationList {
	for i := range opl {
		opl[i].Options.Simulate = simulate
	}
	return opl
}

func (opl OperationList) WithVerbose(verbose bool) OperationList {
	for i := range opl {
		opl[i].Options.Verbose = verbose
	}
	return opl
}

func (opl OperationList) WithIgnoreDirectories(ignoreDirectories bool) OperationList {
	for i := range opl {
		opl[i].Options.IgnoreDirectories = ignoreDirectories
	}
	return opl
}

func (opl OperationList) WithNoIndex(noIndex bool) OperationList {
	for i := range opl {
		opl[i].Options.NoIndex = noIndex
	}
	return opl
}

func (opl OperationList) WithNoExt(noExt bool) OperationList {
	for i := range opl {
		opl[i].Options.NoExt = noExt
	}
	return opl
}

func (opl OperationList) WithNoMkdir(noMkdir bool) OperationList {
	for i := range opl {
		opl[i].Options.NoMkdir = noMkdir
	}
	return opl
}

func (opl OperationList) WithNoMove(noMove bool) OperationList {
	for i := range opl {
		opl[i].Options.NoMove = noMove
	}
	return opl
}

func (opl OperationList) WithSoft(soft bool) OperationList {
	for i := range opl {
		opl[i].Options.Soft = soft
	}
	return opl
}

func (opl OperationList) WithSort(sortOption string) OperationList {
	switch sortOption {
	case SortType.Alphabet:
		sort.Slice(opl, func(i, j int) bool { return opl[i].Input.Abs < opl[j].Input.Abs })
	case SortType.ReverseAlphabet:
		sort.Slice(opl, func(i, j int) bool { return opl[j].Input.Abs < opl[i].Input.Abs })
	case SortType.Date:
		sort.Slice(opl, func(i, j int) bool { return opl[i].Stats.ModTime().Before(opl[j].Stats.ModTime()) })
	case SortType.ReverseDate:
		sort.Slice(opl, func(i, j int) bool { return opl[i].Stats.ModTime().After(opl[j].Stats.ModTime()) })
	case SortType.Size:
		sort.Slice(opl, func(i, j int) bool { return opl[i].Stats.Size() > opl[j].Stats.Size() })
	case SortType.ReverseSize:
		sort.Slice(opl, func(i, j int) bool { return opl[j].Stats.Size() > opl[i].Stats.Size() })
	}
	return opl
}
