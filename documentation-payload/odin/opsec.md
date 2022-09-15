+++
title = "OPSEC"
chapter = false
weight = 10
pre = "<b>1. </b>"
+++

### Post-Exploitation Jobs
All odin commands execute in the context of a go routine or thread. These routines cannot be stopped once started. A "Stop" flag can be sent to them, but if the command itself isn't pausing periodically to check the status of the flag, then it can't be force-exited.

### Agent Compilation
There is currentlyno agent obfuscation.
