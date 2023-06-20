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
cp configs/discord-bot.example.yml configs/discord-bot.yml
vi configs/discord-bot.yml
```

You also can set token via environment variable alternatively.

```bash
export CHATBOT_GPT_DISCORD_TOKEN=token
export CHATBOT_GPT_OPENAI_TOKEN=token
```

## Run

- macOS with Apple Silicon Chip

```bash
./bin/discord-bot/discord-bot_darwin_arm64 --config=configs/discord-bot.yml
```

Otherwise, you need to find the binary file in the `bin` directory,
and specify the configuration file.
