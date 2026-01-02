-- =============================================
-- Portfolio Microservices Database Schema
-- =============================================

-- 1. Auth/User Service Tables
-- =============================================

-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- User Access to Projects
CREATE TABLE IF NOT EXISTS user_project_access (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    project_id INT,
    access_level VARCHAR(50) DEFAULT 'read',
    PRIMARY KEY (user_id, project_id)
);

-- Default roles
INSERT INTO roles (name) VALUES ('admin'), ('user'), ('viewer') ON CONFLICT (name) DO NOTHING;

-- 2. Project/Portfolio Service Tables
-- =============================================

-- Main project table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Skills
CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Project Skills (many-to-many)
CREATE TABLE IF NOT EXISTS project_skills (
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    skill_id INT REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, skill_id)
);

-- Project Tech / stack
CREATE TABLE IF NOT EXISTS project_tech (
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    tech_name VARCHAR(100),
    PRIMARY KEY (project_id, tech_name)
);

-- Project Images
CREATE TABLE IF NOT EXISTS project_images (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    description TEXT,
    uploaded_at TIMESTAMP DEFAULT NOW()
);

-- Project Links
CREATE TABLE IF NOT EXISTS project_links (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    link_url TEXT NOT NULL,
    link_type VARCHAR(50)
);

-- 3. Task/Tracker Service Tables
-- =============================================

-- Main Task
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'Todo',
    priority INT DEFAULT 3,
    assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
    due_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Subtasks
CREATE TABLE IF NOT EXISTS subtasks (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    title VARCHAR(255),
    status VARCHAR(50) DEFAULT 'Todo',
    assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
    due_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Task Comments
CREATE TABLE IF NOT EXISTS task_comments (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Task Attachments
CREATE TABLE IF NOT EXISTS task_attachments (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    file_url TEXT,
    uploaded_at TIMESTAMP DEFAULT NOW()
);

-- Task Tags
CREATE TABLE IF NOT EXISTS task_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE
);

-- Task Tag Mapping
CREATE TABLE IF NOT EXISTS task_tag_mapping (
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INT REFERENCES task_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);

-- 4. Analytics/Stats Service Tables
-- =============================================

-- Project Views
CREATE TABLE IF NOT EXISTS project_views (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    viewed_at TIMESTAMP DEFAULT NOW()
);

-- Task Activity Log
CREATE TABLE IF NOT EXISTS task_activity (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Project Stats
CREATE TABLE IF NOT EXISTS project_stats (
    project_id INT PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    total_tasks INT DEFAULT 0,
    completed_tasks INT DEFAULT 0,
    progress_percent DECIMAL(5,2) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT NOW()
);

-- 5. Media/File Service Tables
-- =============================================

-- Media Files
CREATE TABLE IF NOT EXISTS media_files (
    id SERIAL PRIMARY KEY,
    file_name VARCHAR(255),
    file_url TEXT,
    uploaded_by INT REFERENCES users(id) ON DELETE SET NULL,
    uploaded_at TIMESTAMP DEFAULT NOW(),
    file_type VARCHAR(50)
);

-- =============================================
-- Indexes for better performance
-- =============================================

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_project_views_project_id ON project_views(project_id);
CREATE INDEX IF NOT EXISTS idx_task_activity_task_id ON task_activity(task_id);
CREATE INDEX IF NOT EXISTS idx_media_files_uploaded_by ON media_files(uploaded_by);
