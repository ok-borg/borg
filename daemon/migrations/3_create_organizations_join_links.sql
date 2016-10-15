USE borg;

CREATE TABLE IF NOT EXISTS organization_join_links
(
  id              VARCHAR(36)                         NOT NULL,
  organization_id VARCHAR(36)                         NOT NULL,
  ttl             INTEGER                             NOT NULL,
  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP  NOT NULL,
  created_by      VARCHAR(36)                         NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE organization_join_links
      ADD FOREIGN KEY (organization_id) REFERENCES organizations (id),
      ADD FOREIGN KEY (created_by) REFERENCES users (id);
