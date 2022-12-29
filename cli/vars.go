package main

var Version = "Inserted through ldflags"

// AuthToken is used to authenticate with the Unweave API. It is loaded from the saved
// config file and can be overridden with runtime flags.
var AuthToken = ""

// Path is the path to the current project to run commands on. It is loaded from the saved
// config file and can be overridden with runtime flags.
var Path = ""

// SSHKeyPath is the path to the SSH key to use to connect to a new or existing Session.
var SSHKeyPath = ""
