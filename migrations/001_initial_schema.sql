-- Initial schema for Remnawave Telegram Shop Bot
-- This file contains the initial database schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    language_code VARCHAR(10) DEFAULT 'ru',
    is_blocked BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    balance DECIMAL(10,2) DEFAULT 0.00,
    referral_code VARCHAR(20) UNIQUE,
    referred_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on telegram_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);

-- Create index on referral_code for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_referral_code ON users(referral_code);

-- Create index on referred_by for referral queries
CREATE INDEX IF NOT EXISTS idx_users_referred_by ON users(referred_by);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create servers table
CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on is_active for filtering active servers
CREATE INDEX IF NOT EXISTS idx_servers_is_active ON servers(is_active);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_servers_deleted_at ON servers(deleted_at);

-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    duration INTEGER NOT NULL, -- in days
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on server_id for filtering plans by server
CREATE INDEX IF NOT EXISTS idx_plans_server_id ON plans(server_id);

-- Create index on is_active for filtering active plans
CREATE INDEX IF NOT EXISTS idx_plans_is_active ON plans(is_active);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_plans_deleted_at ON plans(deleted_at);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    server_id INTEGER NOT NULL,
    server_name VARCHAR(255),
    plan_id INTEGER NOT NULL,
    plan_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on user_id for filtering subscriptions by user
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

-- Create index on server_id for filtering subscriptions by server
CREATE INDEX IF NOT EXISTS idx_subscriptions_server_id ON subscriptions(server_id);

-- Create index on plan_id for filtering subscriptions by plan
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);

-- Create index on status for filtering subscriptions by status
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);

-- Create index on expires_at for finding expired subscriptions
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires_at ON subscriptions(expires_at);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_subscriptions_deleted_at ON subscriptions(deleted_at);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'RUB',
    payment_method VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    external_id VARCHAR(255) UNIQUE,
    description VARCHAR(500),
    metadata TEXT, -- JSON with additional data
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Create index on user_id for filtering payments by user
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);

-- Create index on external_id for finding payments by external ID
CREATE INDEX IF NOT EXISTS idx_payments_external_id ON payments(external_id);

-- Create index on status for filtering payments by status
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- Create index on payment_method for filtering payments by method
CREATE INDEX IF NOT EXISTS idx_payments_payment_method ON payments(payment_method);

-- Create index on created_at for sorting payments by date
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_servers_updated_at BEFORE UPDATE ON servers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_plans_updated_at BEFORE UPDATE ON plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default admin user (password should be changed)
INSERT INTO users (telegram_id, username, first_name, last_name, is_admin, referral_code)
VALUES (0, 'admin', 'Admin', 'User', TRUE, 'admin')
ON CONFLICT (telegram_id) DO NOTHING;

-- Insert sample server
INSERT INTO servers (id, name, description, is_active)
VALUES (1, 'Test Server', 'Test server for development', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Insert sample plan
INSERT INTO plans (server_id, name, description, price, duration, is_active)
VALUES (1, 'Basic Plan', 'Basic VPN plan for 30 days', 299.00, 30, TRUE)
ON CONFLICT DO NOTHING;
