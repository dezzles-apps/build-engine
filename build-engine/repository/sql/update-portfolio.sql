UPDATE portfolio_entries
SET title = ?,
    description = ?,
    version = ?,
    url = ?,
    environment = ?,
    public = ?,
    category = ?,
    share_repository = ?,
    repository_url = ?
WHERE project_key = ?