CREATE DATABASE IF NOT EXISTS `zilo_workflow`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `zilo_workflow`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(191) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `nickname` VARCHAR(64) NOT NULL DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflows` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `owner_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(191) NOT NULL,
  `description` VARCHAR(1000) NOT NULL DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `latest_version` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_workflows_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflow_versions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `workflow_id` BIGINT UNSIGNED NOT NULL,
  `version` INT NOT NULL,
  `dsl_json` JSON NOT NULL,
  `is_published` TINYINT NOT NULL DEFAULT 0,
  `created_by` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_workflow_version` (`workflow_id`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflow_executions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `workflow_id` BIGINT UNSIGNED NOT NULL,
  `version_id` BIGINT UNSIGNED NOT NULL,
  `trigger_type` VARCHAR(32) NOT NULL DEFAULT 'manual',
  `status` VARCHAR(32) NOT NULL DEFAULT 'pending',
  `input_json` JSON NULL,
  `output_json` JSON NULL,
  `error_msg` TEXT,
  `started_at` DATETIME NULL,
  `ended_at` DATETIME NULL,
  `duration_ms` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_executions_workflow_started` (`workflow_id`, `started_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflow_execution_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `execution_id` BIGINT UNSIGNED NOT NULL,
  `node_id` VARCHAR(64) NOT NULL,
  `node_type` VARCHAR(64) NOT NULL,
  `status` VARCHAR(32) NOT NULL DEFAULT 'pending',
  `input_json` JSON NULL,
  `output_json` JSON NULL,
  `error_msg` TEXT,
  `started_at` DATETIME NULL,
  `ended_at` DATETIME NULL,
  `duration_ms` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_logs_execution_node` (`execution_id`, `node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

