import mysql from 'mysql2'
let connection
let repository = {
  connect: function(config) {
    return new Promise((resolve, reject) => {
      const conn = mysql.createConnection({
        host: config.host,
        port: config.port,
        user: config.username,
        password: config.password,
        database: config.database
      })

      conn.connect((err) => {
        if (err) {
          console.error('Error connecting to the database:', err)
          reject(err)
        } else {
          connection = conn
          console.log('Connected to the database')
          resolve(repository)
        }
      })
    })
  },
  getConfiguration: function(organisation, repositoryName) {
    return new Promise((resolve, reject) => {
      const query = `
        SELECT r.organisation, r.repository, c.channel
        FROM repositories c
        JOIN repositories r ON c.repository_id = r.repository_id
        WHERE r.organisation = ? AND r.repository = ?
      `
      connection.query(query, [organisation, repositoryName], (err, results) => {
        if (err) {
          console.error('Error fetching configuration:', err)
          reject(err)
        } else {
          if (results.length > 0) {
            resolve(results[0])
          } else {
            resolve(null)
          }
        }
      })
    })
  },

  getBuild: function(organisation, repositoryName, buildNumber) {
    return new Promise((resolve, reject) => {
      const query = `
        SELECT r.organisation, r.repository, b.build_id, b.build_number, b.start_time, b.discord_thread_id
        FROM builds b
        JOIN repositories r ON b.repository_id = r.repository_id
        WHERE r.organisation = ? AND r.repository = ? AND b.build_number = ?
      `
      connection.query(query, [organisation, repositoryName, buildNumber], (err, results) => {
        if (err) {
          console.error('Error fetching build:', err)
          reject(err)
        } else {
          console.log(results)
          if (results.length > 0) {
            resolve(results[0])
          } else {
            resolve(null)
          }
        }
      })
    })
  },

  saveEvent: function(buildId, component, message) {
    return new Promise((resolve, reject) => {
      const query = `
        INSERT INTO build_event (build_id, component, message)
        VALUES (?, ?, ?)
      `
      connection.query(query, [buildId, component, message], (err, results) => {
        if (err) {
          console.error('Error saving event:', err)
          reject(err)
        } else {
          resolve(results)
        }
      })
    })
  },

  setDiscordThreadId: function(organisation, repositoryName, buildNumber, build) {
    return new Promise((resolve, reject) => {
      const query = `
        UPDATE builds b
        JOIN repositories r ON b.repository_id = r.repository_id
        SET b.discord_thread_id = ?
        WHERE r.organisation = ? AND r.repository = ? AND b.build_number = ?
      `
      connection.query(query, [build.discord_thread_id, organisation, repositoryName, buildNumber], (err, results) => {
        if (err) {
          console.error('Error updating discord_thread_id:', err)
          reject(err)
        } else {
          resolve(results)
        }
      })
    })
  }
}


export default repository