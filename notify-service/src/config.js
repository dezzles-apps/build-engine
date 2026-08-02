const env = process.env.ENVIRONMENT || 'non'
const config = {
  non: {
    discord: {
      token: process.env.DISCORD_TOKEN
    },
    database: {
      host: 'dezzles-apps-do-user-35978446-0.h.db.ondigitalocean.com',
      port: 25060,
      username: 'build-engine-non',
      password: process.env.DATABASE_PASSWORD_NON,
      database: 'build-engine'
    }
  },
  prod: {
    discord: {
      token: process.env.DISCORD_TOKEN
    },
    database: {
      host: 'dezzles-apps-do-user-35978446-0.h.db.ondigitalocean.com',
      port: 25060,
      username: 'build-engine-prd',
      password: process.env.DATABASE_PASSWORD_PROD,
      database: 'build-engine'
    }
  }
}

export default config[env]