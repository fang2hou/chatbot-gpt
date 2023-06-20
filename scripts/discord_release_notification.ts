import { readFile } from "fs/promises";
import {
  APIEmbedField,
  EmbedAuthorOptions,
  EmbedBuilder,
  RestOrArray,
  WebhookClient,
} from "discord.js";
import { Configuration, OpenAIApi } from "openai";

// env returns the value of the environment variable with the given name.
let env = (envName: string): string => {
  let value = process.env[envName];
  if (value == null || value.length === 0) {
    throw new Error(`Environment variable ${envName} is not set`);
  }
  return value;
};

// getArtifactUrl returns the URL of the artifact with the given name.
const getArtifactUrl = (prefix: string, artifactName: string): string => {
  return `${artifactsUrlBase}/${prefix}_${artifactName}`;
};

// artifactSuffixes is a map of artifact names to their suffixes.
const artifactSuffixes: { [key: string]: string } = {
  "macOS 英特尔芯片": "darwin_amd64",
  "macOS Apple 芯片": "darwin_arm64",
  "Linux AMD 64 位": "linux_amd64",
  "Linux ARM 64 位": "linux_arm64",
  "Linux ARMv6 芯片": "linux_armv6",
  "Linux ARMv7 芯片": "linux_armv7",
  "Windows AMD 64 位": "windows_amd64.exe",
  "Windows ARM 64 位": "windows_arm64.exe",
};

const REPOSITORY: string = env("GITHUB_REPOSITORY");
const TAG: string = env("GITHUB_REF_NAME");
const GITHUB_URL: string = env("GITHUB_SERVER_URL");
const DISCORD_WEBHOOK_URL: string = env("DISCORD_WEBHOOK_URL");
const OPENAI_API_KEY: string = env("OPENAI_API_KEY");

const projectName: string = REPOSITORY.split("/")[1];
const releaseUrl: string = `${GITHUB_URL}/${REPOSITORY}/releases/tag/${TAG}`;
const artifactsUrlBase: string = `${GITHUB_URL}/${REPOSITORY}/releases/download/${TAG}`;

async function main() {
  let changes: string = await readFile("./changes.md", "utf8");
  let releaseNote: string = "";

  changes.split("\n").forEach((line: string) => {
    if (line.startsWith("build")) {
      releaseNote += `• :tools: ${line.trim()}\n`;
    } else if (line.startsWith("feat")) {
      releaseNote += `• :star2: ${line.trim()}\n`;
    } else if (line.startsWith("fix")) {
      releaseNote += `• :bug: ${line.trim()}\n`;
    } else if (line.startsWith("refactor")) {
      releaseNote += `• :recycle: ${line.trim()}\n`;
    } else if (line.startsWith("docs")) {
      releaseNote += `• :books: ${line.trim()}\n`;
    } else if (line.startsWith("style")) {
      releaseNote += `• :art: ${line.trim()}\n`;
    } else if (line.startsWith("perf")) {
      releaseNote += `• :zap: ${line.trim()}\n`;
    } else if (line.startsWith("test")) {
      releaseNote += `• :white_check_mark: ${line.trim()}\n`;
    } else if (line.startsWith("chore")) {
      releaseNote += `• :wrench: ${line.trim()}\n`;
    } else if (line.startsWith("revert")) {
      releaseNote += `• :rewind: ${line.trim()}\n`;
    } else if (line.startsWith("ci")) {
      releaseNote += `• :robot: ${line.trim()}\n`;
    } else if (line.startsWith("deps")) {
      releaseNote += `• :package: ${line.trim()}\n`;
    } else if (line.startsWith("misc")) {
      releaseNote += `• :label: ${line.trim()}\n`;
    } else if (line.length === 0) {
    } else {
      releaseNote += `${line}\n`;
    }
  });

  // Translate the release note to Chinese
  let openaiClient = new OpenAIApi(
    new Configuration({ apiKey: OPENAI_API_KEY })
  );
  let translationResp = await openaiClient.createChatCompletion({
    model: "gpt-3.5-turbo-0301",
    messages: [
      {
        role: "user",
        content:
          "You are a translator for translating the software changelog to chinese," +
          "your answer should not include anything outside the list." +
          "Do not translate the word matched ':.*:', just keep it as it is.",
      },
      {
        role: "user",
        content: releaseNote,
      },
    ],
  });

  if (translationResp.data.choices[0].message != null) {
    releaseNote = translationResp.data.choices[0].message.content;
    releaseNote += "\n\n:beginner: 以上内容由 OpenAI GPT-3.5 Turbo 生成";
  }

  const fields: RestOrArray<APIEmbedField> = [];

  fields.push(
    { name: "**变更列表**", value: releaseNote, inline: false },
    { name: "\u200B", value: "\u200B", inline: false }
  );

  for (const [key, value] of Object.entries(artifactSuffixes)) {
    fields.push({
      name: `:low_brightness: ${key}`,
      value: `:small_blue_diamond: [下载](${getArtifactUrl(
        "discord-bot",
        value
      )})`,
      inline: true,
    });
  }

  let author: EmbedAuthorOptions = {
    name: "GitHub Actions",
    iconURL: "https://avatars.githubusercontent.com/in/15368",
    url: releaseUrl,
  };

  const embed: EmbedBuilder = new EmbedBuilder()
    .setTitle(`:loudspeaker: ${projectName} ${TAG} 版本发布了！`)
    .setAuthor(author)
    .setURL(releaseUrl)
    .addFields(fields)
    .setColor(0x1baf9c);

  const webhookClient: WebhookClient = new WebhookClient({
    url: DISCORD_WEBHOOK_URL,
  });

  await webhookClient.send({ embeds: [embed] }).catch((error) => {
    console.error(error);
  });
  webhookClient.destroy();
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
