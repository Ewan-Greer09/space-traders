-- Create the users table
CREATE TABLE players (
    id INT AUTO_INCREMENT PRIMARY KEY, user_uid VARCHAR(255) UNIQUE, username VARCHAR(100), password VARCHAR(100), email VARCHAR(100), created_at TIMESTAMP
);

-- create the api_keys table
CREATE TABLE api_keys (
    id INT AUTO_INCREMENT PRIMARY KEY, api_key VARCHAR(1000), u_id VARCHAR(255), FOREIGN KEY (u_id) REFERENCES players (user_uid)
);

-- create the agents table
CREATE TABLE agents (
    id INT AUTO_INCREMENT PRIMARY KEY, accountId VARCHAR(255), symbol VARCHAR(255), head_quarters VARCHAR(255), credits INT, starting_faction VARCHAR(255), ship_count INT
);

-- create the systems TABLE
CREATE TABLE systems (
    id INT AUTO_INCREMENT PRIMARY KEY, symbol VARCHAR(255) UNIQUE, sector_symbol VARCHAR(255), type VARCHAR(255), x INT, y INT
);

-- create the waypoints table
CREATE TABLE waypoints (
    id INT AUTO_INCREMENT PRIMARY KEY, system_symbol VARCHAR(255), symbol VARCHAR(255) UNIQUE, type VARCHAR(255), x INT, y INT, orbits VARCHAR(255), FOREIGN KEY (system_symbol) REFERENCES systems (symbol)
);

-- create table orbitals
CREATE TABLE orbitals (
    id INT AUTO_INCREMENT PRIMARY KEY, symbol VARCHAR(255), waypoint_symbol VARCHAR(255) UNIQUE, FOREIGN KEY (waypoint_symbol) REFERENCES waypoints (symbol)
);

-- create the factions table
CREATE TABLE factions (
    id INT AUTO_INCREMENT PRIMARY KEY, symbol VARCHAR(255), waypoint_symbol VARCHAR(255), FOREIGN KEY (waypoint_symbol) REFERENCES waypoints (symbol)
);
