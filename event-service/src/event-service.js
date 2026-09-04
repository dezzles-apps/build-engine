import Discord from './discord.js'
import Repository from './repository.js'

let discordClient
let eventRepository

let eventService = {
  init: function(discord, r) {
    discordClient = discord
    eventRepository = r
  },
  event: async function({ organisation, repository, buildNumber, source, sourceBuildId, message, component, ref, buildUrl }) {
    try {
      const [build, config] = await Promise.all([
        eventRepository.getBuild(source, sourceBuildId),
        eventRepository.getConfiguration(organisation, repository)
      ])
      if (!build) {
        throw new Error(`Build not found for ${organisation}/${repository} #${buildNumber} (source: ${source}, sourceBuildId: ${sourceBuildId})`)
      }
      if (!config) {
        throw new Error(`Configuration not found for ${organisation}/${repository}`)
      }
      const channel = config.channel ? config.channel : 'dezzles-apps'
      console.log(build.discord_thread_id)
      if (!build.discord_thread_id) {
        let threadId = await discordClient.createThread(channel, organisation, repository, buildNumber, source, sourceBuildId, ref, buildUrl)
        build.discord_thread_id = threadId.id
        console.log(`Created new thread for ${organisation}/${repository} #${buildNumber}: ${build.discord_thread_id}`)
        await eventRepository.setDiscordThreadId(source, sourceBuildId, build)
      } else {
        console.log(`Thread already exists for ${organisation}/${repository} #${buildNumber}: ${build.discord_thread_id}`)
      }
      await eventRepository.saveEvent(build.build_id, component, message)
      await discordClient.sendMessage(channel, build.discord_thread_id, message)
    } catch (error) {
      console.error('Error in eventService.event:', error)
      throw error
    }
  }
}

export default eventService