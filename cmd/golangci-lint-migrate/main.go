package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/bombsimon/wsl/v5"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string `yaml:"version"`
}

type V1 struct {
	Settings Settings `yaml:"linters-settings"`
}

type V2 struct {
	Linters Linters `yaml:"linters"`
}

type Linters struct {
	Settings Settings `yaml:"settings"`
}

type Settings struct {
	WSL *WSL `yaml:"wsl"`
}

type NewConfig struct {
	Linters struct {
		Settings struct {
			WSL WSLV5 `yaml:"wsl"`
		} `yaml:"settings"`
	} `yaml:"linters"`
}

// StrictAppend                 is replaced with CheckAppend
// AllowAssignAndCall           is replaced with CheckAssignExpr (inverse)
// AllowAssignAndAnything       is replaced with CheckAssign
// AllowMultilineAssign         is deprecated and always allowed in v5
// AllowTrailingComment         is deprecated and always allowed in v5
// AllowCuddleDeclarations      is replaced with CheckDecl
// AllowSeparatedLeadingComment is deprecated and always allowed in v5
// ForceErrCuddling             is replaced with CheckErr
// ForceShortDeclCuddling       is replaced with CheckAssignExclusive
// ForceCaseTrailingWhitespace  is deprecated and replaced with CaseMaxLines
// AllowCuddlingWithCalls       is deprecated and not needed in v5
// AllowCuddleWithRHS           is deprecated and not needed in v5
// ErrorVariableNames           is deprecated and not needed in v5
type WSL struct {
	StrictAppend                 bool `yaml:"strict-append"`
	AllowAssignAndCall           bool `yaml:"allow-assign-and-call"`
	AllowAssignAndAnything       bool `yaml:"allow-assign-and-anything"`
	AllowMultilineAssign         bool `yaml:"allow-multiline-assign"`
	AllowTrailingComment         bool `yaml:"allow-trailing-comment"`
	AllowCuddleDeclarations      bool `yaml:"allow-cuddle-declarations"`
	AllowSeparatedLeadingComment bool `yaml:"allow-separated-leading-comment"`
	ForceErrCuddling             bool `yaml:"force-err-cuddling"`
	ForceShortDeclCuddling       bool `yaml:"force-short-decl-cuddling"`
	ForceCaseTrailingWhitespace  int  `yaml:"force-case-trailing-whitespace"`

	AllowCuddlingWithCalls []string `yaml:"allow-cuddle-with-calls"`
	AllowCuddleWithRHS     []string `yaml:"allow-cuddle-with-rhs"`
	ErrorVariableNames     []string `yaml:"error-variable-names"`
}

type WSLV5 struct {
	AllowFirstInBlock bool     `yaml:"allow-first-in-block"`
	AllowWholeBlock   bool     `yaml:"allow-whole-block"`
	BranchMaxLines    int      `yaml:"branch-max-lines"`
	CaseMaxLines      int      `yaml:"case-max-lines"`
	Enable            []string `yaml:"enable"`
	Disable           []string `yaml:"disable"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("expected exactly one argument which should be your `.golangci.yml` file")
	}

	v1cfg := getWslConfig(os.Args[1])
	if v1cfg == nil {
		log.Println("failed to get wsl configuration from config")
		log.Println("ensure you passed a configuration with existing `wsl` configuration")
		log.Fatalln("if you didn't use any configuration you can keep leaving it empty to use the defaults")
	}

	v5cfg := WSLV5{
		AllowFirstInBlock: true,
		AllowWholeBlock:   false,
		BranchMaxLines:    2,
		CaseMaxLines:      v1cfg.ForceCaseTrailingWhitespace,
	}

	if !v1cfg.StrictAppend {
		v5cfg.Disable = append(v5cfg.Disable, wsl.CheckAppend.String())
	}

	if v1cfg.AllowAssignAndAnything {
		v5cfg.Disable = append(v5cfg.Disable, wsl.CheckAssign.String())
	}

	if v1cfg.AllowMultilineAssign {
		log.Println("`allow-multiline-assign` is deprecated and always allowed in >= v5")
	}

	if v1cfg.AllowTrailingComment {
		log.Println("`allow-trailing-comment` is deprecated and always allowed in >= v5")
	}

	if v1cfg.AllowSeparatedLeadingComment {
		log.Println("`allow-separated-leading-comment` is deprecated and always allowed in >= v5")
	}

	if v1cfg.AllowCuddleDeclarations {
		v5cfg.Disable = append(v5cfg.Disable, wsl.CheckDecl.String())
	}

	if v1cfg.ForceErrCuddling {
		v5cfg.Enable = append(v5cfg.Enable, wsl.CheckErr.String())
	}

	if v1cfg.ForceShortDeclCuddling {
		v5cfg.Enable = append(v5cfg.Enable, wsl.CheckAssignExclusive.String())
	}

	if !v1cfg.AllowAssignAndCall {
		v5cfg.Enable = append(v5cfg.Enable, wsl.CheckAssignExpr.String())
	}

	slices.Sort(v5cfg.Enable)
	slices.Sort(v5cfg.Disable)

	cfg := NewConfig{}
	cfg.Linters.Settings.WSL = v5cfg

	buf := bytes.NewBuffer([]byte{})
	e := yaml.NewEncoder(buf)
	e.SetIndent(2)

	if err := e.Encode(&cfg); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("")
	fmt.Println("These settings are the closest you can get in the new version of `wsl`")
	fmt.Println("Potential deprecations are logged above")
	fmt.Println("")
	fmt.Println("See https://github.com/bombsimon/wsl for more details")
	fmt.Println("")

	fmt.Println(buf.String())
}

func getWslConfig(filename string) *WSL {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("%v", err)
	}

	switch cfg.Version {
	case "":
		var cfg V1
		if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
			log.Fatalf("%v", err)
		}

		return cfg.Settings.WSL
	case "2":
		var cfg V2
		if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
			log.Fatalf("%v", err)
		}

		return cfg.Linters.Settings.WSL
	default:
		log.Fatalf("invalid version '%s'", cfg.Version)
	}

	return nil
}
