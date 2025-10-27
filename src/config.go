package main

import (
	cfg "grpcd/config"
)

// Shim: keep a reference to the shared config package to avoid duplicate
// declarations in this tree. Other files should import cfg and use
// cfg.AppConfig directly.
var _ = cfg.AppConfig
 
