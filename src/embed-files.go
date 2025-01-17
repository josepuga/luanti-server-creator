package main

import (
    "embed"
)
//TODO: Copiarlos con os.CopyFS
//go:embed embed/others/*
var othersFS embed.FS

//go:embed embed/scripts/*
var scriptsFS embed.FS


