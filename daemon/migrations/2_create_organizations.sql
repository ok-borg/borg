USE borg;

CREATE TABLE IF NOT EXISTS organizations
(
  id            VARCHAR(36)                         NOT NULL,
  name          VARCHAR(512)                        NOT NULL,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP  NOT NULL,
  updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP  NOT NULL,
  created_by    VARCHAR(36)                         NOT NULL,
  updated_by    VARCHAR(36)                         NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS user_organizations
(
  id              VARCHAR(36)                         NOT NULL,
  user_id         VARCHAR(36)                         NOT NULL,
  organization_id VARCHAR(36)                         NOT NULL,
  is_admin        TINYINT DEFAULT 0                   NOT NULL,
  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP  NOT NULL,
  updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP  NOT NULL,
  created_by      VARCHAR(36)                         NOT NULL,
  updated_by      VARCHAR(36)                         NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE organizations
      ADD FOREIGN KEY (created_by) REFERENCES users (id),
      ADD FOREIGN KEY (updated_by) REFERENCES users (id);


ALTER TABLE user_organizations
      ADD FOREIGN KEY (user_id) REFERENCES users (id),
      ADD FOREIGN KEY (organization_id) REFERENCES organizations (id),
      ADD FOREIGN KEY (created_by) REFERENCES users (id),
      ADD FOREIGN KEY (updated_by) REFERENCES users (id);
