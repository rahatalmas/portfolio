-- Initial schema. Statements are separated by a semicolon at end of line.

CREATE TABLE IF NOT EXISTS settings (
  `key`   VARCHAR(64) NOT NULL,
  `value` TEXT        NOT NULL,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS education (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  degree     VARCHAR(255) NOT NULL,
  school     VARCHAR(255) NOT NULL,
  years      VARCHAR(64)  NOT NULL DEFAULT '',
  note       VARCHAR(500) NOT NULL DEFAULT '',
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_education_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS skill_groups (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name       VARCHAR(100) NOT NULL,
  items      TEXT         NOT NULL,
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_skill_groups_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profiles (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  group_name VARCHAR(100) NOT NULL,
  label      VARCHAR(100) NOT NULL,
  href       VARCHAR(500) NOT NULL,
  icon       VARCHAR(32)  NOT NULL,
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_profiles_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS experience (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title      VARCHAR(255) NOT NULL,
  company    VARCHAR(255) NOT NULL,
  period     VARCHAR(64)  NOT NULL,
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_experience_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS projects (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  slug       VARCHAR(100) NOT NULL,
  title      VARCHAR(255) NOT NULL,
  summary    TEXT         NOT NULL,
  tags       VARCHAR(255) NOT NULL DEFAULT '',
  stack      TEXT         NOT NULL,
  featured   TINYINT(1)   NOT NULL DEFAULT 0,
  github_url VARCHAR(500) NOT NULL DEFAULT '',
  live_url   VARCHAR(500) NOT NULL DEFAULT '',
  paper_url  VARCHAR(500) NOT NULL DEFAULT '',
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uq_projects_slug (slug),
  KEY idx_projects_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS papers (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title      VARCHAR(500) NOT NULL,
  authors    TEXT         NOT NULL,
  venue      VARCHAR(255) NOT NULL,
  year       INT          NOT NULL,
  status     VARCHAR(20)  NOT NULL DEFAULT 'published',
  pdf_url    VARCHAR(500) NOT NULL DEFAULT '',
  code_url   VARCHAR(500) NOT NULL DEFAULT '',
  poster_url VARCHAR(500) NOT NULL DEFAULT '',
  slides_url VARCHAR(500) NOT NULL DEFAULT '',
  sort_order INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_papers_sort (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS messages (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name       VARCHAR(120) NOT NULL,
  email      VARCHAR(254) NOT NULL,
  body       TEXT         NOT NULL,
  is_read    TINYINT(1)   NOT NULL DEFAULT 0,
  created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_messages_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS admin_users (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username      VARCHAR(64)  NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_admin_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
