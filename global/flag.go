package global

import (
	"github.com/selefra/selefra-utils/pkg/pointer"
	"sync"
)

var WORKSPACE = pointer.ToStringPointer(".")
var LOGINTOKEN = ""
var ORGNAME = ""
var CMD = ""
var STAG = ""
var LOGLEVEL = "debug"
var levelMap = map[string]bool{
	"trace":   true,
	"debug":   true,
	"info":    true,
	"warning": true,
	"error":   true,
	"fatal":   true,
}

var o sync.Once

func ChangeLevel(level string) {
	if levelMap[level] {
		o.Do(func() {
			LOGLEVEL = level
		})
	}
}

const PkgBasePath = "ghcr.io/selefra/postgre_"
const PkgTag = ":latest"

var SERVER = "main-api.selefra.io"
