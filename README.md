# Freyja Purple Team Agent

<p align="center">
  <img alt="Freyja Logo" src="documentation-payload/freyja/freyja.svg" height="30%" width="30%">
</p>

Freyja is a Golang, Purple Team agent that compiles into Windows, Linux and macOS x64 executables.  It is a very stripped down verion of the [Poseidon](https://github.com/MythicAgents/poseidon) payload from [@xorrior](https://github.com/xorrior), [@djhohnstein](https://github.com/djhohnstein), and [@its-a-feature](https://github.com/its-a-feature)

This Freyja instance supports Mythic 3.0.0 and will be updated as necessary.
It does not support Mythic 2.3 and lower.

The agent has `mythic_payloadtype_container==0.1.8` PyPi package installed and reports to Mythic as version "12".

Freyja uses Red Canary's [Atomic Red Team (ART)](https://github.com/redcanaryco/atomic-red-team) "executor" concept.  Freyja uses these executors to run commands on the victim host machines.  Freyja will integrate with the upcoming Mythic Purple Team eXecution Framework (PTXF) to run ART atomics (test cases), custom atomics via yaml files, and atomic chains within fully customizable Purple Team campagins.  Stay tuned for the PTXF!

Current Executors:
- Windows powershell
- Windows cmd
- macOS zsh
- Linux/macOS bash
- Linux/macOS sh

## How to install an agent in this format within Mythic

When it's time for you to test out your install or for another user to install your agent, it's pretty simple. Within Mythic you can run the `mythic-cli` binary to install this in one of three ways:

* `sudo ./mythic-cli install github https://github.com/user/repo` to install the main branch
* `sudo ./mythic-cli install github https://github.com/user/repo branchname` to install a specific branch of that repo
* `sudo ./mythic-cli install folder /path/to/local/folder/cloned/from/github` to install from an already cloned down version of an agent repo

Now, you might be wondering _when_ should you or a user do this to properly add your agent to their Mythic instance. There's no wrong answer here, just depends on your preference. The three options are:

* Mythic is already up and going, then you can run the install script and just direct that agent's containers to start (i.e. `sudo ./mythic-cli start agentName` and if that agent has its own special C2 containers, you'll need to start them too via `sudo ./mythic-cli start c2profileName`).
* Mythic is already up and going, but you want to minimize your steps, you can just install the agent and run `sudo ./mythic-cli start`. That script will first _stop_ all of your containers, then start everything back up again. This will also bring in the new agent you just installed.
* Mythic isn't running, you can install the script and just run `sudo ./mythic-cli start`. 

## Documentation

The Freyja documentation source code can be found in the `documenation-payload/freyja` directory.
View the rendered documentation by clicking on **Docs -> Agent Documentation** in the upper right-hand corner of the Mythic
interface. 

## Building Outside of Mythic

If you want to build outside of Mythic, you can use the `Makefile` included in the project's `agent_code` directory. You will need to modify the variables at the top of the Makefile to match the C2 profile information you want to build into your agent. To get all the pieces you need (like UUID and AES key), you need to build the agent within Mythic (or at least kick off an unsuccessful build), then copy that information. To find the information you need, simply go to the Payloads page and click the blue info icon. You'll see the UUID, encryption key, and any other information you need for building to put into your Makefile.
