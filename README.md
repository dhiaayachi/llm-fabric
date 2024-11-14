![llmfabric workflow](https://github.com/dhiaayachi/llm-fabric/actions/workflows/build.yaml/badge.svg)
# llm-fabric

`llm-fabric` is a Go-based framework designed to allow multiple AI agents to cooperate to solve problems using a
provided strategy it also provide a unified interface for AI tasks submission and an agent discovery module. 
The framework supports multiple LLM providers, such as OpenAI, Ollama, and Anthropic, 
and allows easy addition of new providers.

## Features

- **Multi-LLM Support**: Unified interface for multiple LLM providers (e.g., OpenAI, Ollama, Anthropic).
- **Extensible Interface**: Easily add new LLM providers, discoverer, strategy...
- **gRPC API**: Expose LLM capabilities through a gRPC service.
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

### Environment Variables

Set up environment variables to configure API keys for different providers:

- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`



## Development

### Running Unit Tests

```bash
go test ./... -cover
```

### Contributing

TBD