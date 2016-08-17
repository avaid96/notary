package main

import "github.com/docker/docker-credential-helpers/osxkeychain"

const defaultCredentialsStore = "osxkeychain"
var helper  = osxkeychain.Osxkeychain{}