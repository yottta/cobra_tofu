package commands

import "fmt"

// =============================
// Moved here from OpenTofu repo strictly to for testing purposes.
// Once we will move towards adoption of cobra, this kind of custom
// types will be less and less necessary.
// =============================

// rawFlags is a flag.Value implementation that just appends raw flag
// names and values to a slice.
type rawFlags struct {
	flagName string
	items    *[]rawFlag
}

func newRawFlags(flagName string) rawFlags {
	var items []rawFlag
	return rawFlags{
		flagName: flagName,
		items:    &items,
	}
}

func (f rawFlags) Empty() bool {
	if f.items == nil {
		return true
	}
	return len(*f.items) == 0
}

func (f rawFlags) AllItems() []rawFlag {
	if f.items == nil {
		return nil
	}
	return *f.items
}

func (f rawFlags) Alias(flagName string) rawFlags {
	return rawFlags{
		flagName: flagName,
		items:    f.items,
	}
}

func (f rawFlags) String() string {
	return ""
}

func (f rawFlags) Set(str string) error {
	*f.items = append(*f.items, rawFlag{
		Name:  f.flagName,
		Value: str,
	})
	return nil
}

// This method satisfies pflag.Value
func (f rawFlags) Type() string {
	return "stringSlice"
}

type rawFlag struct {
	Name  string
	Value string
}

func (f rawFlag) String() string {
	return fmt.Sprintf("%s=%q", f.Name, f.Value)
}
