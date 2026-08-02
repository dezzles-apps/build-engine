import Discord from './discord.js'
import Koa from 'koa'
import Router from '@koa/router'
import bodyParser from '@koa/bodyparser'
import Repository from './repository.js'
import EventService from './event-service.js'
import Config from './config.js'
console.log('Beginning connection to Discord and Database...')
await Promise.all([Discord.connect(Config.discord), Repository.connect(Config.database)])
console.log('Connected to Discord and Database.')
EventService.init(Discord, Repository)
const app = new Koa()
const router = new Router()
app.use(bodyParser())
router.post('/api/v1/notify', async (ctx) => {
  const { organisation, repository, buildNumber, message, component, ref } = ctx.request.body;
  if (!component) {
    ctx.request.body.component = 'default'
  }
  if (!organisation || !repository || !buildNumber || !message || !ref) {
    ctx.status = 400;
    ctx.body = { error: 'Missing required fields' };
    return;
  }

  try {
    await EventService.event(ctx.request.body);
    ctx.status = 200;
    ctx.body = { success: true };
  } catch (error) {
    console.error('Error sending message to Discord:', error);
    ctx.status = 500;
    ctx.body = { error: 'Failed to send message to Discord' };
  }
})

app.use(router.routes())

app.listen(process.env.SERVER_PORT || 3000, () => {
  console.log(`Server running on port ${process.env.SERVER_PORT || 3000}`);
});