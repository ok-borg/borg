USE borg;

CREATE TABLE IF NOT EXISTS users
(
  id                  VARCHAR(36)                                                     NOT NULL,
  username            VARCHAR(512)                                                    NOT NULL,
  email               VARCHAR(512)                                                    NOT NULL,
  github_id           VARCHAR(512)                                                    NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


