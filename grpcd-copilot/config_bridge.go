package main

import cfg "grpcd/config"

// AppConfigProxy exposes the shared configuration for unit tests that still
// reference the old global symbol. This keeps backward compatibility while
// delegating all real work to the cfg package.
var AppConfig = &cfg.AppConfig

func configInit() {
	cfg.Init()
}

func LoadConfigDefault() {
	cfg.LoadConfigDefault()
}
