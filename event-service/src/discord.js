import { Client, Events, GatewayIntentBits } from 'discord.js';

let persistedClient

function title(organisation, repository, buildNumber) {
  return `${organisation}/${repository} - #${buildNumber}`
}

function newThreadMessage(event) {
  const t = `${event.organisation}/${event.repository} - #${event.buildNumber}`
  return `${t}
Ref: ${event.ref}
[Github Actions Run](${event.buildUrl})`
}

let discordClient = {
  connect: async function(config) {
    return new Promise((resolve, reject) => {
      const client = new Client({ intents: [GatewayIntentBits.Guilds] });
      async function onConnect(connectedClient) {
        console.log(`Ready! Logged in as ${client.user.tag}`);
        persistedClient = connectedClient
        resolve(discordClient);
      }
      client.once(Events.ClientReady, onConnect);
      client.login(config.token)
    })
  },
  createThread: async function(channelName, event) {
    const threadName = `${event.organisation}/${event.repository} #${event.buildNumber}`;
    console.log(`Creating thread: ${threadName} in channel: ${channelName}`);
    const channel = persistedClient.channels.cache.find(c => c.name === channelName);
    const message = await channel.send(newThreadMessage(event));
    console.log('messageId', message.id);
    return channel.threads.create({
      name: title(event.organisation, event.repository, event.buildNumber),
      autoArchiveDuration: 60,
      startMessage: message.id,
      type: 11, // 11 is the type for a public thread
    }).then(thread => {
      console.log(`Created thread with ID: ${thread.id}`)
      return thread
    }).catch(error => {
      console.error(`Failed to create thread: ${error.message}`)
      throw error
    })
  },
  sendMessage: async function(channelName, threadId, message) {
    const channel = persistedClient.channels.cache.find(c => c.name === channelName);
    const thread = channel.threads.cache.get(threadId);
    if (!thread) {
      throw new Error(`Thread with ID ${threadId} not found in channel ${channelName}`);
    }
    return thread.send(message);
  }
}

export default discordClient
