USE borg;

CREATE TABLE IF NOT EXISTS users
(
  id                  VARCHAR(36)                                                     NOT NULL,
  login               VARCHAR(512)                                                    NOT NULL,
  name                VARCHAR(512)                                                    NOT NULL,
  email               VARCHAR(512)                                                    NOT NULL,
  avatar_url          VARCHAR(512)                                                    NOT NULL,
-- account type will be a string with the name of oauth provider (github, google, etc) 
  account_type        VARCHAR(512)                                                    NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS github_users
(
  id                  VARCHAR(36)                                                     NOT NULL,
  github_id           VARCHAR(36)                                                     NOT NULL,
  borg_user_id        VARCHAR(36)                                                     NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE github_users
      ADD FOREIGN KEY (borg_user_id) REFERENCES users (id);

CREATE TABLE IF NOT EXISTS access_tokens
(
  id                  VARCHAR(36)                                                     NOT NULL,
  token               VARCHAR(512)                                                    NOT NULL,
  user_id             VARCHAR(36)                                                     NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE access_tokens
      ADD FOREIGN KEY (user_id) REFERENCES users (id);
