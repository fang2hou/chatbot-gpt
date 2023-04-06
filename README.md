# chatbot-gpt

Chat bots powered by [OpenAI](https://openai.com/)'s GPT API.

## Support Platforms

- [x] Discord
- [ ] Telegram
- [ ] Slack

## Build

To build this project, you need to install [Go](https://golang.org/) and [Make](https://www.gnu.org/software/make/).

### Build for current platform

```bash
make -f makefile build-all
```

### Cross compile for all platforms

```bash
make -f makefile cross-build-all
```

## Prepare Configurations

```bash
cp configs/discord-bot.example.json configs/discord-bot.json
vi configs/discord-bot.json
```

You also can set token via environment variable alternatively.

```bash
export CHATBOTS_GPT_DISCORD_TOKEN=token
export CHATBOTS_GPT_OPENAI_TOKEN=token
```

## Run

- macOS with Apple Silicon Chip

```bash
./bin/discord-bot/discord-bot_darwin_arm64 --config=configs/discord-bot.yaml
```

Otherwise, you need to find the binary file in the `bin` directory, and specify the configuration file.
