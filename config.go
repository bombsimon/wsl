package wsl

import (
	"errors"
	"fmt"
	"strings"
)

// CheckSet is a set of checks to run.
type CheckSet map[CheckType]struct{}

// CheckType is a type that represents a checker to run.
type CheckType int

// Each checker is represented by a CheckType that is used to enable or disable
// the check.
// A check can either be of a specific built-in keyword or custom checks.
const (
	CheckInvalid CheckType = iota
	CheckAssign
	CheckBreak
	CheckContinue
	CheckDecl
	CheckDefer
	CheckExpr
	CheckFor
	CheckGo
	CheckIf
	CheckIncDec
	CheckLabel
	CheckRange
	CheckReturn
	CheckSelect
	CheckSend
	CheckSwitch
	CheckTypeSwitch

	// Assign exclusive only allows assignments of either new variables or
	// re-assignment of existing ones, e.g.
	//
	// a := 1
	// b := 2
	//
	// a = 1
	// b = 2
	CheckAssignExclusive
	// Check append only allows assignments of `append` to be cuddled with other
	// assignments if it's a variable used in the append statement, e.g.
	//
	// a := 1
	// x = append(x, a)
	CheckAppend
	// Force error checking to follow immediately after an error variable is
	// assigned, e.g.
	//
	// _, err := someFn()
	// if err != nil {
	//     panic(err)
	// }
	CheckErr
	CheckLeadingWhitespace
	CheckTrailingWhitespace
)

/*
Configuration migration

- StrictAppend                     | Deprecate. Replaced with `CheckAppend`
- AllowAssignAndCallCuddle         | TBD deprecate. Implemented, not configurable
- AllowAssignAndAnythingCuddle     | Deprecated. Replaced with `CheckAssign`
- AllowMultiLineAssignCuddle       | Deprecate.
- ForceCaseTrailingWhitespaceLimit | Done.
- AllowTrailingComment             | TBD deprecate. Should be seen same as leading (allowed)
- AllowSeparatedLeadingComment     | Deprecate. Always allowed.
- AllowCuddleDeclaration           | Deprecate. Use `CheckDecl` instead.
- AllowCuddleWithCalls             | TBD deprecate. Should not be needed. Was added to support mutex unlocking
- AllowCuddleWithRHS               | TBD deprecate. Should not be needed. Not clear why separate from above
- ForceCuddleErrCheckAndAssign     | Deprecate. Replaced with `CheckErr`
- ErrorVariableNames               | Deprecate. We're now looking if the variable implements the error interface
- ForceExclusiveShortDeclarations  | Done.
- IncludeGenerated                 | Done.
*/
type Configuration struct {
	AllowFirstInBlock bool
	AllowWholeBlock   bool
	CaseMaxLines      int
	ReturnMaxLines    int
	Checks            CheckSet
}

func NewConfig() *Configuration {
	return &Configuration{
		AllowFirstInBlock: true,
		AllowWholeBlock:   false,
		CaseMaxLines:      0,
		ReturnMaxLines:    2,
		Checks:            DefaultChecks(),
	}
}

func (c *Configuration) Update(
	enableAll bool,
	disableAll bool,
	enable []string,
	disable []string,
) error {
	if enableAll && disableAll {
		return errors.New("can't use both `enable-all` and `disable-all`")
	}

	if enableAll {
		c.Checks = AllChecks()
	}

	if disableAll {
		c.Checks = NoChecks()
	}

	for _, s := range enable {
		check, err := CheckFromString(s)
		if err != nil {
			return fmt.Errorf("invalid check '%s'", s)
		}

		c.Checks.Add(check)
	}

	for _, s := range disable {
		check, err := CheckFromString(s)
		if err != nil {
			return fmt.Errorf("invalid check '%s'", s)
		}

		c.Checks.Remove(check)
	}

	return nil
}

func DefaultChecks() CheckSet {
	return CheckSet{
		CheckAssign:             {},
		CheckAppend:             {},
		CheckBreak:              {},
		CheckContinue:           {},
		CheckDecl:               {},
		CheckDefer:              {},
		CheckExpr:               {},
		CheckFor:                {},
		CheckGo:                 {},
		CheckIf:                 {},
		CheckIncDec:             {},
		CheckLabel:              {},
		CheckLeadingWhitespace:  {},
		CheckTrailingWhitespace: {},
		CheckRange:              {},
		CheckReturn:             {},
		CheckSelect:             {},
		CheckSwitch:             {},
		CheckTypeSwitch:         {},
	}
}

func AllChecks() CheckSet {
	c := DefaultChecks()
	c.Add(CheckAssignExclusive)
	c.Add(CheckErr)
	c.Add(CheckSend)

	return c
}

func NoChecks() CheckSet {
	return CheckSet{}
}

func (c CheckSet) Add(check CheckType) {
	c[check] = struct{}{}
}

func (c CheckSet) Remove(check CheckType) {
	delete(c, check)
}

func CheckFromString(s string) (CheckType, error) {
	switch strings.ToLower(s) {
	case "assign":
		return CheckAssign, nil
	case "break":
		return CheckBreak, nil
	case "continue":
		return CheckContinue, nil
	case "decl":
		return CheckDecl, nil
	case "defer":
		return CheckDefer, nil
	case "expr":
		return CheckExpr, nil
	case "for":
		return CheckFor, nil
	case "go":
		return CheckGo, nil
	case "if":
		return CheckIf, nil
	case "incdec":
		return CheckIncDec, nil
	case "label":
		return CheckLabel, nil
	case "range":
		return CheckRange, nil
	case "return":
		return CheckReturn, nil
	case "select":
		return CheckSelect, nil
	case "send":
		return CheckSend, nil
	case "switch":
		return CheckSwitch, nil
	case "type-switch":
		return CheckTypeSwitch, nil

	case "append":
		return CheckAppend, nil
	case "assign-exclusive":
		return CheckAssignExclusive, nil
	case "err":
		return CheckErr, nil
	case "leading-whitespace":
		return CheckLeadingWhitespace, nil
	case "trailing-whitespace":
		return CheckTrailingWhitespace, nil
	default:
		return CheckInvalid, fmt.Errorf("invalid check '%s'", s)
	}
}
