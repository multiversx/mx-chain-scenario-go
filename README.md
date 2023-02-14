# MultiversX blockchain scenarios: The Go framework

## Overview

Blockchain scenarios are interactions with the blockchain, real or imagined.

They help developers write tests and are able to document and replay interactions with smart contracts.

The format is described here: https://docs.multiversx.com/developers/scenario-reference/overview


## Scope

This Go framework deals with reading, writing, and controlling scenario runners.

Scenario runners are the routines that do something with these scenarios. Think of them as closures that receive the scenario steps.

The main example for such a runner can be found in the VM, here: https://github.com/multiversx/mx-chain-vm-go/tree/master/scenarioexec

However, more such runners are conceivable.

To implement such a runner, create an object that implements interface `ScenarioExecutor`.

## Alternate implementation

There is an equivalent Rust implementation here: https://github.com/multiversx/mx-sdk-rs/tree/master/sdk/scenario-format

The Go implementation (this one) is older and generally tends to be better featured, altough they should be up-to-date with one another now.


