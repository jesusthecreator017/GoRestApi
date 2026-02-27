-- Admin user. Password: Admin123!  (change after first login)
INSERT INTO users (email, password_hash, name, permissions)
VALUES ('admin@fswithgo.com', '$2a$10$.FxzphTRqpQp.x8lxNT6Yec6hvcoPgP9jaRsnw6JReuiRF/scaARe', 'Admin', 7)
ON CONFLICT (email) DO UPDATE SET permissions = 7;
