-- Создание таблицы для записей СЗИ
CREATE TABLE IF NOT EXISTS szi_records (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    szi_type VARCHAR(255) NOT NULL,
    cert_number VARCHAR(255) NOT NULL,
    cert_issue_date DATE NOT NULL,
    cert_expiry_date DATE NOT NULL,
    location VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'Актуально',
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    manufacturer VARCHAR(255),
    software_version VARCHAR(255),
    contact_person VARCHAR(255),
    documentation_link TEXT
);

-- Создание таблицы для пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Создание таблицы для прав доступа
CREATE TABLE IF NOT EXISTS access_permissions (
    id SERIAL PRIMARY KEY,
    owner_user_id INTEGER NOT NULL,
    guest_user_id INTEGER NOT NULL,
    record_id INTEGER NOT NULL,
    permission VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Создание таблицы для ролей пользователей
CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    role_name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы для уведомлений
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    read_status BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP WITH TIME ZONE
);

-- Создание таблицы для запросов доступа
CREATE TABLE IF NOT EXISTS access_requests (
    id SERIAL PRIMARY KEY,
    requester_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    record_id INTEGER NOT NULL,
    permission VARCHAR(10) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE,
    processed_by INTEGER
);