-- Migration: Create users table for high-concurrency auth system
-- Supports 100,000+ concurrent users

-- Drop existing table if it exists
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) UNIQUE,
    email VARCHAR(255) UNIQUE,
    username VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    birthday DATE,
    latest_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure at least one of email, username, or phone_number is provided
    CONSTRAINT check_contact_info CHECK (
        email IS NOT NULL OR username IS NOT NULL OR phone_number IS NOT NULL
    )
);

-- Indexes for high-performance queries
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_latest_login ON users(latest_login);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update latest_login
CREATE OR REPLACE FUNCTION update_latest_login(user_id BIGINT)
RETURNS VOID AS $$
BEGIN
    UPDATE users 
    SET latest_login = CURRENT_TIMESTAMP 
    WHERE id = user_id;
END;
$$ language 'plpgsql';
