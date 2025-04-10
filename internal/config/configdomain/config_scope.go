package configdomain

import "strings"

// defines the type of Git configuration used
type ConfigScope string

const (
	ConfigScopeGlobal ConfigScope = "global"
	ConfigScopeLocal  ConfigScope = "local"
)

func (self ConfigScope) String() string {
	return string(self)
}

func ParseConfigScope(text string) ConfigScope {
	switch strings.TrimSpace(text) {
	case "local", "":
		return ConfigScopeLocal
	case "global":
		return ConfigScopeGlobal
	default:
		panic("unknown locality: " + text)
	}
}
