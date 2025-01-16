![llmfabric workflow](https://github.com/dhiaayachi/llm-fabric/actions/workflows/build.yaml/badge.svg)
# llm-fabric

`llm-fabric` is a Go-based framework designed to allow multiple AI agents to cooperate to solve problems using a
provided strategy it also provide a unified interface for AI tasks submission and an agent discovery module. 
The framework supports multiple LLM providers, such as OpenAI, Ollama, and Anthropic, 
and allows easy addition of new providers.

## Features

- **Multi-LLM Support**: Unified interface for multiple LLM providers (e.g., OpenAI, Ollama, Anthropic).
- **Extensible Interface**: Easily add new LLM providers, discoverer, strategy...
- **Configurable**: Customizable capabilities, tools.

## Usage

see the examples folder, for some examples on how to use `llm-fabric`

## Running examples

Clone the repository and install dependencies:

```bash
git clone https://github.com/dhiaayachi/llm-fabric.git
cd llm-fabric/examples/dispatcher
export OPENAI_API_KEY=<key>
docker compose build
docker compose up
```

## Configuration

Configuration can be done by setting up specific parameters when initializing each LLM client. For example, `NewOllama` and `NewGPT` allow setting capabilities, tools, and logging preferences.
See the [examples](https://github.com/dhiaayachi/llm-fabric/tree/e594fa250646d915baf59f40f7c2ff4ea7ca392a/examples) for more details

### Environment Variables

Set up environment variables to configure API keys for different providers:

- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`

## Implementation

llm-fabric is a framework composed of multiple components. Most components are wrapped in an interface to allow 
multiple implementations and extensibility.

### Components

#### Fabric

This is the core component that would instantiate a new fabric. It include an `Agent` or a `Client`. 

- The `Agent`purpose is to present an Llm to the Fabric and allow it to be available to solve tasks. 
A special type of agents are `Dispatchers`, they allow to dispatch tasks to other agents using `Strategies`.
- The `Client` purpose is to allow an application to discover agents and will forward any task submitted by 
the application to a dispatcher agent to dispatch it to a capable agent to solve it.

#### Llm

This is a thin wrapper around a given llm API and provide a common API to interact with those LLMs. 
It's implemented for `GPT`, `Ollama`, `Anthropic` for the time being.
A new implementation could be added by implementing 
the [llm interface](https://github.com/dhiaayachi/llm-fabric/blob/main/llm/llm.go#L58-L58)

#### Discoverer

This component is responsible for discovering other LLMs added to the fabric. The current implementation leverage 
[hashicorp/serf](github.com/hashicorp/serf) library. A new implementation need to implement 
the [dicoverer interface](https://github.com/dhiaayachi/llm-fabric/blob/main/discoverer/discoverer.go#L11-L11)

#### Agent

This is responsible for implementing the communication interface between all the agents (using gRPC) and 
allow agents to instantiate gRPC servers and clients 

#### Strategy

A strategy is the logic that a `Dispatcher` agent will use to dispatch a task to agents to solve it.
A strategy example would be to [Select the most capable agent with the least cost](https://github.com/dhiaayachi/llm-fabric/blob/main/examples/dispatcher/dispatcher_ollama/capability_dispatcher.go#L81-L81).

### Running Unit Tests

```bash
go test ./... -cover
```

### Contributing

TBD