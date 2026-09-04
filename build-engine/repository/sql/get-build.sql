SELECT 
  b.source_build_id,
  r.organisation,
  r.repository,
  b.build_number,
  b.start_time,
  b.discord_thread_id 
FROM 
  builds b
JOIN repositories r ON b.repository_id = r.repository_id
WHERE
  b.source = ? AND b.source_build_id = ?