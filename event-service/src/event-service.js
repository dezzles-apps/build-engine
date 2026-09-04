import Discord from './discord.js'
import Repository from './repository.js'

let discordClient
let eventRepository

let eventService = {
  init: function(discord, r) {
    discordClient = discord
    eventRepository = r
  },
  event: async function(event) {
    try {
      const [build, config] = await Promise.all([
        eventRepository.getBuild(event.source, event.sourceBuildId),
        eventRepository.getConfiguration(event.organisation, event.repository)
      ])
      if (!build) {
        throw new Error(`Build not found for ${event.organisation}/${event.repository} #${event.buildNumber} (source: ${event.source}, sourceBuildId: ${event.sourceBuildId})`)
      }
      if (!config) {
        throw new Error(`Configuration not found for ${event.organisation}/${event.repository}`)
      }
      const channel = config.channel ? config.channel : 'dezzles-apps'
      console.log(build.discord_thread_id)
      if (!build.discord_thread_id) {
        let threadId = await discordClient.createThread(channel, event)
        build.discord_thread_id = threadId.id
        console.log(`Created new thread for ${event.organisation}/${event.repository} #${event.buildNumber}: ${build.discord_thread_id}`)
        await eventRepository.setDiscordThreadId(event.source, event.sourceBuildId, build)
      } else {
        console.log(`Thread already exists for ${event.organisation}/${event.repository} #${event.buildNumber}: ${build.discord_thread_id}`)
      }
      await eventRepository.saveEvent(build.build_id, event.component, event.message)
      await discordClient.sendMessage(channel, build.discord_thread_id, event.message)
    } catch (error) {
      console.error('Error in eventService.event:', error)
      throw error
    }
  }
}

export default eventService