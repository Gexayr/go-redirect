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
    request_body JSON,
    response_status INT,
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
    redirect_type ENUM('site1', 'site2', 'site3') NOT NULL,
    redirect_status INT NOT NULL,
    redirect_timestamp DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_log_id) REFERENCES request_logs(id),
    INDEX idx_redirect_type (redirect_type),
    INDEX idx_redirect_timestamp (redirect_timestamp)
);

-- Request Parameters Table (depends on request_logs)
CREATE TABLE IF NOT EXISTS request_parameters (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_log_id BIGINT NOT NULL,
    parameter_name VARCHAR(255) NOT NULL,
    parameter_value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_log_id) REFERENCES request_logs(id),
    INDEX idx_request_log_id (request_log_id)
);

-- Error Logs Table (depends on request_logs)
CREATE TABLE IF NOT EXISTS error_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_log_id BIGINT,
    error_message TEXT NOT NULL,
    error_type VARCHAR(100) NOT NULL,
    stack_trace TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_log_id) REFERENCES request_logs(id),
    INDEX idx_error_type (error_type)
);

-- Insert sample redirect mappings
INSERT INTO redirect_mappings (hash, redirect_url) VALUES
('123456', 'http://site1.com'),
('789012', 'http://site2.com'),
('345678', 'http://site3.com'),
('any-random-hash', 'http://site1.com/special-offer'); 