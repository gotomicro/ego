package eflag

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	flagset *FlagSet
)

func init() {
	flagset = &FlagSet{
		FlagSet: flag.CommandLine,
		flags:   defaultFlags,
		actions: make(map[string]func(string, *FlagSet)),
	}
}

// Flag ...
type (
	// Flag defines application flag.
	Flag interface {
		Apply(*FlagSet)
	}

	// FlagSet wraps a set of Flags.
	FlagSet struct {
		*flag.FlagSet
		flags   []Flag
		actions map[string]func(string, *FlagSet)
	}
)

// SetFlagSet 设置flagSet
func SetFlagSet(fs *FlagSet) {
	flagset = fs
}

// NewFlagSet new flagSet
func NewFlagSet(flagCommand *flag.FlagSet, flags ...Flag) *FlagSet {
	return &FlagSet{
		FlagSet: flagCommand,
		flags:   flags,
		actions: make(map[string]func(string, *FlagSet)),
	}
}

// Register ...
func Register(fs ...Flag) {
	flagset.Register(fs...)
}

// Register ...
func (fs *FlagSet) Register(flags ...Flag) {
	fs.flags = append(fs.flags, flags...)
}

// With adds flags to the flagset.
func With(fs ...Flag) { flagset.With(fs...) }

// With adds flags to provided flagset.
func (fs *FlagSet) With(flags ...Flag) {
	fs.flags = append(fs.flags, flags...)
}

// Parse parses the flagset.
func Parse() error {
	return flagset.Parse()
}

// ParseWithArgs parses the flagset with given args.
func ParseWithArgs(arguments []string) error {
	return flagset.ParseWithArgs(arguments)
}

// Lookup lookup flag value by name
// priority: flag > env > default
func (fs *FlagSet) Lookup(name string) *flag.Flag {
	return fs.FlagSet.Lookup(name)
}

// Parse parses provided flagset.
func (fs *FlagSet) Parse() error {
	return fs.ParseWithArgs(os.Args[1:])
}

// ParseWithArgs parses provided flagset with given args.
func (fs *FlagSet) ParseWithArgs(arguments []string) error {
	if fs.Parsed() {
		return nil
	}
	for _, f := range fs.flags {
		f.Apply(fs)
	}

	// 解析命令行参数
	if err := fs.FlagSet.Parse(arguments); err != nil {
		return err
	}

	// 遍历所欲flagset数据
	fs.FlagSet.Visit(func(f *flag.Flag) {
		if action, ok := fs.actions[f.Name]; ok && action != nil {
			action(f.Name, fs)
		}
	})
	return nil
}

// BoolE parses bool flag of the flagset with error returned.
func BoolE(name string) (bool, error) { return flagset.BoolE(name) }

// BoolE parses bool flag of provided flagset with error returned.
func (fs *FlagSet) BoolE(name string) (bool, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseBool(flag.Value.String())
	}

	return false, fmt.Errorf("undefined flag name: %s", name)
}

// Bool parses bool flag of the flagset.
func Bool(name string) bool { return flagset.Bool(name) }

// Bool parses bool flag of provided flagset.
func (fs *FlagSet) Bool(name string) bool {
	ret, _ := fs.BoolE(name)
	return ret
}

// StringE parses string flag of the flagset with error returned.
func StringE(name string) (string, error) { return flagset.StringE(name) }

// StringE parses string flag of provided flagset with error returned.
func (fs *FlagSet) StringE(name string) (string, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return flag.Value.String(), nil
	}

	return "", fmt.Errorf("undefined flag name: %s", name)
}

// String parses string flag of the flagset.
func String(name string) string { return flagset.String(name) }

// String parses string flag of provided flagset.
func (fs *FlagSet) String(name string) string {
	ret, _ := fs.StringE(name)
	return ret
}

// IntE parses int flag of the flagset with error returned.
func IntE(name string) (int64, error) { return flagset.IntE(name) }

// IntE parses int flag of provided flagset with error returned.
func (fs *FlagSet) IntE(name string) (int64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseInt(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

// Int parses int flag of the flagset.
func Int(name string) int64 { return flagset.Int(name) }

// Int parses int flag of provided flagset.
func (fs *FlagSet) Int(name string) int64 {
	ret, _ := fs.IntE(name)
	return ret
}

// UintE parses int flag of the flagset with error returned.
func UintE(name string) (uint64, error) { return flagset.UintE(name) }

// UintE parses int flag of provided flagset with error returned.
func (fs *FlagSet) UintE(name string) (uint64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseUint(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

// Uint parses int flag of the flagset.
func Uint(name string) uint64 { return flagset.Uint(name) }

// Uint parses int flag of provided flagset.
func (fs *FlagSet) Uint(name string) uint64 {
	ret, _ := fs.UintE(name)
	return ret
}

// Float64E parses int flag of the flagset with error returned.
func Float64E(name string) (float64, error) { return flagset.Float64E(name) }

// Float64E parses int flag of provided flagset with error returned.
func (fs *FlagSet) Float64E(name string) (float64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseFloat(flag.Value.String(), 64)
	}

	return 0.0, fmt.Errorf("undefined flag name: %s", name)
}

// Float64 parses int flag of the flagset.
func Float64(name string) float64 { return flagset.Float64(name) }

// Float64 parses int flag of provided flagset.
func (fs *FlagSet) Float64(name string) float64 {
	ret, _ := fs.Float64E(name)
	return ret
}
