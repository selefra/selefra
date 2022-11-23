package global

import (
	"github.com/selefra/selefra-utils/pkg/pointer"
)

var WORKSPACE = pointer.ToStringPointer(".")
var LOGINTOKEN = ""
var ORGNAME = ""
var CMD = ""
var STAG = ""

const PkgBasePath = "ghcr.io/selefra/postgre_"
const PkgTag = ":latest"

var SERVER = "main-api.selefra.io"
