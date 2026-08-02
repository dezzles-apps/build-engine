
CREATE TABLE repositories (
  repository_id int AUTO_INCREMENT PRIMARY KEY,
  organisation varchar(255) NOT NULL,
  repository varchar(255) NOT NULL,
  UNIQUE KEY unique_repository (organisation, repository),
  channel varchar(255) NOT NULL
);

CREATE TABLE builds (
  build_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  repository_id int,
  build_number int,
  start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (repository_id) REFERENCES repositories(repository_id)
);

ALTER TABLE builds ADD COLUMN discord_thread_id BIGINT;

CREATE TABLE build_event(
  build_event_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  build_id BIGINT,
  time DATETIME DEFAULT CURRENT_TIMESTAMP,
  message VARCHAR(255) NOT NULL,
  component VARCHAR(255) NOT NULL,
  FOREIGN KEY (build_id) REFERENCES builds(build_id)
);

ALTER TABLE build-engine.builds ADD COLUMN ref varchar(255) NOT NULL;

CREATE TABLE valid_tokens(
  token_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  token varchar(255) NOT NULL UNIQUE
);
