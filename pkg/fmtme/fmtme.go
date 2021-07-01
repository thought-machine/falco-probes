package fmtme

import "golang.org/x/xerrors"

// FmtMe is a temp var used to help test linting/fmting
var FmtMe = "temporary package to test plz fmt-all command"

// ImportMe is a temp var used to test go mod tidy works through plz fmt-all
var ImportMe = xerrors.As
