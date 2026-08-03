DROP TABLE IF EXISTS demo;
DROP TABLE IF EXISTS priorities;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS todos;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(30) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS statuses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  emoji VARCHAR(5) NOT NULL,
  name VARCHAR(15) UNIQUE NOT NULL
);

INSERT INTO statuses 
(id, emoji, name)
VALUES
(1, '📝', 'Planning'),
(2, '🟢', 'Active'),
(3, '📈', 'In Progress'),
(4, '❌', 'Cancelled'),
(5, '🗃️', 'Archive'),
(6, '🗑️', 'Trash');

CREATE TABLE IF NOT EXISTS priorities (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(15) UNIQUE NOT NULL,
  emoji VARCHAR(5) NOT NULL
);

INSERT INTO priorities 
(id, emoji, name)
VALUES
(1, '🔴', 'Important'),
(2, '🟠', 'Highest'),
(3, '🟡', 'High'),
(4, '🟢', 'Medium'),
(5, '🔵', 'Low'),
(6, '⚫', 'Lowest');

CREATE TABLE IF NOT EXISTS todos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task VARCHAR(300) NOT NULL,
  due TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id INTEGER NOT NULL,
  status_id INTEGER DEFAULT 1,
  priority_id INTEGER DEFAULT 3,
  
  CONSTRAINT fk_todos_users FOREIGN KEY (user_id) REFERENCES users(id)
  ON UPDATE CASCADE ON DELETE CASCADE,
  
  CONSTRAINT fk_todos_statuses FOREIGN KEY (status_id) REFERENCES statuses(id)
  ON UPDATE CASCADE ON DELETE RESTRICT,
  
  CONSTRAINT fk_todos_priorities FOREIGN KEY (priority_id) REFERENCES priorities(id)
  ON UPDATE CASCADE ON DELETE RESTRICT
);

PRAGMA foreign_keys = ON;
