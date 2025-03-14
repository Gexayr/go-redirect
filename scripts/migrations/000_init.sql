-- Create database if not exists
CREATE DATABASE IF NOT EXISTS platform_db;
USE platform_db;

-- Clients Table
CREATE TABLE IF NOT EXISTS clients (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
);

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
    client_id BIGINT NOT NULL,
    hash VARCHAR(6) NOT NULL UNIQUE,
    redirect_url VARCHAR(255) NOT NULL,
    redirect_url_black VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id),
    INDEX idx_hash (hash),
    INDEX idx_client_id (client_id)
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

-- Insert sample client (password: test123)
INSERT INTO clients (username, password_hash, email) VALUES
('testuser', '$2a$10$G00Q5kIlzJFA7jSxyBrdJek.lZTgcuwOcj94AxGIfp2DZK57Sc55e', 'test@example.com');

-- Insert sample redirect mappings
INSERT INTO redirect_mappings (client_id, hash, redirect_url, redirect_url_black) VALUES
(1, '123456', 'http://youtube.com', 'http://youtube.com'),
(1, '789012', 'http://goodwin.am', 'http://goodwin.am'),
(1, '345678', 'http://google.com', 'http://google.com'),
(1, 'random', 'https://chatgpt.com', 'https://chatgpt.com');
