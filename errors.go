package main

import "errors"

var (
	ErrInvalidInputType  = errors.New("invalid input type")
	ErrInvalidOutputType = errors.New("invalid output type")
	ErrConfigFileMissing = errors.New("configuration file is missing")
)
