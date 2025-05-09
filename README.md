# gollm

gollm is a Go library which provides a simple interface for interacting with any LLM in Go.

## Usage

```bash
go get github.com/dickeyy/gollm
```

```go
model, err := gollm.InitializeModel("gpt-4o-mini", "your-api-key")
if err != nil { panic(err) }

res, err := model.Chat(gollm.ChatStructure{
    Messages: []gollm.ChatMessage{
        {
            Role: "user",
            Content: "Hello, how are you?"
        },
    },
})
if err != nil { panic(err) }

fmt.Println(res.Text)
```

## Supported Models

- OpenAI:
  - GPT-4o-mini (`gpt-4o-mini`)
  - GPT-4o (`gpt-4o`)
  - GPT-3.5-turbo (`gpt-3.5-turbo`)

_More models coming soon!_

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
