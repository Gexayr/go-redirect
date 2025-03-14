-- Create database if not exists
CREATE DATABASE IF NOT EXISTS platform_db;
USE platform_db;

-- Request Logs Table (base table with no foreign keys)
CREATE TABLE IF NOT EXISTS request_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    request_url TEXT NOT NULL,
    request_method VARCHAR(10) NOT NULL,
    request_headers JSON,
    processing_status ENUM('pending', 'processed', 'failed') DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_timestamp (timestamp),
    INDEX idx_processing_status (processing_status)
);

-- Redirect Mappings Table (for URL mappings)
CREATE TABLE IF NOT EXISTS redirect_mappings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    hash VARCHAR(255) NOT NULL UNIQUE,
    redirect_url VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_hash (hash)
);

-- Redirect History Table (depends on request_logs)
CREATE TABLE IF NOT EXISTS redirect_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_log_id BIGINT NOT NULL,
    original_url TEXT NOT NULL,
    redirect_url TEXT NOT NULL,
    redirect_type VARCHAR(50) NOT NULL,
    redirect_status INT NOT NULL,
    redirect_timestamp DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_log_id) REFERENCES request_logs(id),
    INDEX idx_redirect_type (redirect_type),
    INDEX idx_redirect_timestamp (redirect_timestamp)
);

-- Insert sample redirect mappings
INSERT INTO redirect_mappings (hash, redirect_url) VALUES
('123456', 'http://youtube.com'),
('789012', 'http://goodwin.am'),
('345678', 'http://google.com'),
('any-random-hash', 'https://chatgpt.com');