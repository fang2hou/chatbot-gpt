import {readFile} from 'fs';
import {APIEmbedField, EmbedAuthorOptions, EmbedBuilder, RestOrArray, WebhookClient} from 'discord.js';


const REPOSITORY: string = process.env.GITHUB_REPOSITORY || '';
const TAG: string = process.env.GITHUB_REF_NAME || '';
const GITHUB_URL: string = process.env.GITHUB_SERVER_URL || '';
const DISCORD_WEBHOOK_URL: string = process.env.DISCORD_WEBHOOK_URL || '';

if (REPOSITORY.length === 0) {
    console.error('GITHUB_REPOSITORY is not set');
    process.exit(1);
}

if (TAG.length === 0) {
    console.error('GITHUB_REF_NAME is not set');
    process.exit(1);
}

if (GITHUB_URL.length === 0) {
    console.error('GITHUB_SERVER_URL is not set');
    process.exit(1);
}

if (DISCORD_WEBHOOK_URL.length === 0) {
    console.error('DISCORD_WEBHOOK_URL is not set');
    process.exit(1);
}

const projectName: string = REPOSITORY.split('/')[1];
const releaseUrl: string = `${GITHUB_URL}/${REPOSITORY}/releases/tag/${TAG}`;
const artifactsUrlBase: string = `${GITHUB_URL}/${REPOSITORY}/releases/download/${TAG}`;

readFile('./changes.md', 'utf8', (err: NodeJS.ErrnoException | null, data: string) => {
    if (err) {
        console.error(`Error reading changes.md: ${err}`);
        return;
    }

    const fields: RestOrArray<APIEmbedField> = [];

    let releaseNote: string = '';
    // Create the formatted message content
    data.split('\n').forEach((line: string) => {
        if (line.startsWith('build')) {
            releaseNote += `• :tools: ${line.trim()}\n`;
        } else if (line.startsWith('feat')) {
            releaseNote += `• :star2: ${line.trim()}\n`;
        } else if (line.startsWith('fix')) {
            releaseNote += `• :bug: ${line.trim()}\n`;
        } else if (line.startsWith('refactor')) {
            releaseNote += `• :recycle: ${line.trim()}\n`;
        } else if (line.startsWith('docs')) {
            releaseNote += `• :books: ${line.trim()}\n`;
        } else if (line.startsWith('style')) {
            releaseNote += `• :art: ${line.trim()}\n`;
        } else if (line.startsWith('perf')) {
            releaseNote += `• :zap: ${line.trim()}\n`;
        } else if (line.startsWith('test')) {
            releaseNote += `• :white_check_mark: ${line.trim()}\n`;
        } else if (line.startsWith('chore')) {
            releaseNote += `• :wrench: ${line.trim()}\n`;
        } else if (line.startsWith('revert')) {
            releaseNote += `• :rewind: ${line.trim()}\n`;
        } else if (line.startsWith('ci')) {
            releaseNote += `• :robot: ${line.trim()}\n`;
        } else if (line.startsWith('deps')) {
            releaseNote += `• :package: ${line.trim()}\n`;
        } else if (line.startsWith('misc')) {
            releaseNote += `• :label: ${line.trim()}\n`;
        } else if (line.length === 0) {
        } else {
            releaseNote += `${line}\n`;
        }
    });

    fields.push(
        {name: '**变更列表**', value: releaseNote, inline: false},
        {name: '\u200B', value: '\u200B', inline: false}
    );

    const getArtifactUrl = (prefix: string, artifactName: string): string => {
        return `${artifactsUrlBase}/${prefix}_${artifactName}`;
    };

    const artifactSuffixes: { [key: string]: string } = {
        'macOS 英特尔芯片': 'darwin_amd64',
        'macOS Apple 芯片': 'darwin_arm64',
        'Linux AMD 64 位': 'linux_amd64',
        'Linux ARM 64 位': 'linux_arm64',
        'Linux ARMv6 芯片': 'linux_armv6',
        'Linux ARMv7 芯片': 'linux_armv7',
        'Windows AMD 64 位': 'windows_amd64.exe',
        'Windows ARM 64 位': 'windows_arm64.exe',
    };

    for (const [key, value] of Object.entries(artifactSuffixes)) {
        fields.push({
            name: `:low_brightness: ${key}`,
            value: `[:small_blue_diamond: 下载](${getArtifactUrl("discord-bot", value)})`,
            inline: true
        });
    }

    let author: EmbedAuthorOptions = {
        name: 'GitHub Actions',
        iconURL: 'https://avatars.githubusercontent.com/in/15368',
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

    webhookClient
        .send({embeds: [embed]})
        .then(() => {
            console.log('Discord notification sent successfully!');
        })
        .catch((error: Error) => {
            console.error('Failed to send Discord notification:', error);
        })
        .finally(() => {
            webhookClient.destroy();
        });
});