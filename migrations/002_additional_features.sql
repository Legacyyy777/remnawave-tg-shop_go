-- Additional features migration for Remnawave Telegram Shop Bot
-- This file contains additional tables for promo codes, notifications, and activity logs

-- Create promo_codes table
CREATE TABLE IF NOT EXISTS promo_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(20) DEFAULT 'bonus_days',
    value DECIMAL(10,2) NOT NULL,
    max_uses INTEGER DEFAULT 0,
    used_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    valid_until TIMESTAMP WITH TIME ZONE,
    description VARCHAR(500),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for promo_codes
CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON promo_codes(code);
CREATE INDEX IF NOT EXISTS idx_promo_codes_type ON promo_codes(type);
CREATE INDEX IF NOT EXISTS idx_promo_codes_is_active ON promo_codes(is_active);
CREATE INDEX IF NOT EXISTS idx_promo_codes_valid_from ON promo_codes(valid_from);
CREATE INDEX IF NOT EXISTS idx_promo_codes_valid_until ON promo_codes(valid_until);
CREATE INDEX IF NOT EXISTS idx_promo_codes_created_by ON promo_codes(created_by);
CREATE INDEX IF NOT EXISTS idx_promo_codes_deleted_at ON promo_codes(deleted_at);

-- Create promo_code_usages table
CREATE TABLE IF NOT EXISTS promo_code_usages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    promo_code_id UUID NOT NULL REFERENCES promo_codes(id),
    user_id UUID NOT NULL REFERENCES users(id),
    used_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for promo_code_usages
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_promo_code_id ON promo_code_usages(promo_code_id);
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_user_id ON promo_code_usages(user_id);
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_used_at ON promo_code_usages(used_at);

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_is_sent ON notifications(is_sent);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Create activity_logs table
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    data TEXT, -- JSON with additional data
    ip_address VARCHAR(45), -- IPv4 or IPv6
    user_agent VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for activity_logs
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at);

-- Create triggers for new tables
CREATE TRIGGER update_promo_codes_updated_at BEFORE UPDATE ON promo_codes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add trial_used column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS trial_used BOOLEAN DEFAULT FALSE;

-- Create index for trial_used
CREATE INDEX IF NOT EXISTS idx_users_trial_used ON users(trial_used);

-- Add referral_bonus_earned column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_bonus_earned DECIMAL(10,2) DEFAULT 0.00;

-- Create index for referral_bonus_earned
CREATE INDEX IF NOT EXISTS idx_users_referral_bonus_earned ON users(referral_bonus_earned);

-- Add last_activity_at column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Create index for last_activity_at
CREATE INDEX IF NOT EXISTS idx_users_last_activity_at ON users(last_activity_at);

-- Add language preference to users table (if not exists)
ALTER TABLE users ADD COLUMN IF NOT EXISTS language_preference VARCHAR(10) DEFAULT 'ru';

-- Create index for language_preference
CREATE INDEX IF NOT EXISTS idx_users_language_preference ON users(language_preference);

-- Add subscription status indexes for better performance
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_status ON subscriptions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status_expires ON subscriptions(status, expires_at);

-- Add payment status indexes for better performance
CREATE INDEX IF NOT EXISTS idx_payments_user_status ON payments(user_id, status);
CREATE INDEX IF NOT EXISTS idx_payments_method_status ON payments(payment_method, status);

-- Create function to clean up old activity logs
CREATE OR REPLACE FUNCTION cleanup_old_activity_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM activity_logs 
    WHERE created_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- Create function to clean up old notifications
CREATE OR REPLACE FUNCTION cleanup_old_notifications()
RETURNS void AS $$
BEGIN
    DELETE FROM notifications 
    WHERE created_at < NOW() - INTERVAL '30 days' 
    AND is_read = TRUE;
END;
$$ LANGUAGE plpgsql;

-- Create function to get user statistics
CREATE OR REPLACE FUNCTION get_user_stats(user_id_param UUID)
RETURNS TABLE (
    total_subscriptions BIGINT,
    active_subscriptions BIGINT,
    total_payments DECIMAL(10,2),
    total_spent DECIMAL(10,2),
    referral_count BIGINT,
    last_activity TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(s.id) as total_subscriptions,
        COUNT(CASE WHEN s.status = 'active' AND s.expires_at > NOW() THEN 1 END) as active_subscriptions,
        COALESCE(SUM(p.amount), 0) as total_payments,
        COALESCE(SUM(CASE WHEN p.status = 'completed' THEN p.amount ELSE 0 END), 0) as total_spent,
        COUNT(r.id) as referral_count,
        u.last_activity_at as last_activity
    FROM users u
    LEFT JOIN subscriptions s ON u.id = s.user_id
    LEFT JOIN payments p ON u.id = p.user_id
    LEFT JOIN users r ON u.id = r.referred_by
    WHERE u.id = user_id_param
    GROUP BY u.id, u.last_activity_at;
END;
$$ LANGUAGE plpgsql;

-- Create function to get admin statistics
CREATE OR REPLACE FUNCTION get_admin_stats()
RETURNS TABLE (
    total_users BIGINT,
    active_users BIGINT,
    blocked_users BIGINT,
    total_subscriptions BIGINT,
    active_subscriptions BIGINT,
    total_revenue DECIMAL(10,2),
    today_revenue DECIMAL(10,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(u.id) as total_users,
        COUNT(CASE WHEN u.last_activity_at > NOW() - INTERVAL '7 days' THEN 1 END) as active_users,
        COUNT(CASE WHEN u.is_blocked = TRUE THEN 1 END) as blocked_users,
        COUNT(s.id) as total_subscriptions,
        COUNT(CASE WHEN s.status = 'active' AND s.expires_at > NOW() THEN 1 END) as active_subscriptions,
        COALESCE(SUM(CASE WHEN p.status = 'completed' THEN p.amount ELSE 0 END), 0) as total_revenue,
        COALESCE(SUM(CASE WHEN p.status = 'completed' AND p.completed_at::date = CURRENT_DATE THEN p.amount ELSE 0 END), 0) as today_revenue
    FROM users u
    LEFT JOIN subscriptions s ON u.id = s.user_id
    LEFT JOIN payments p ON u.id = p.user_id;
END;
$$ LANGUAGE plpgsql;
