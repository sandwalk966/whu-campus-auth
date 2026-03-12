INSERT INTO sys_user (username,password,nickname,email,phone,status,created_at,updated_at) VALUES
('user001','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser001','user001@test.com','13800000001',1,NOW(),NOW()),
('user002','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser002','user002@test.com','13800000002',1,NOW(),NOW()),
('user003','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser003','user003@test.com','13800000003',1,NOW(),NOW()),
('user004','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser004','user004@test.com','13800000004',1,NOW(),NOW()),
('user005','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser005','user005@test.com','13800000005',1,NOW(),NOW()),
('user006','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser006','user006@test.com','13800000006',1,NOW(),NOW()),
('user007','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser007','user007@test.com','13800000007',1,NOW(),NOW()),
('user008','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser008','user008@test.com','13800000008',1,NOW(),NOW()),
('user009','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser009','user009@test.com','13800000009',1,NOW(),NOW()),
('user010','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser010','user010@test.com','13800000010',1,NOW(),NOW()),
('user011','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser011','user011@test.com','13800000011',1,NOW(),NOW()),
('user012','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser012','user012@test.com','13800000012',1,NOW(),NOW()),
('user013','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser013','user013@test.com','13800000013',1,NOW(),NOW()),
('user014','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser014','user014@test.com','13800000014',1,NOW(),NOW()),
('user015','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser015','user015@test.com','13800000015',1,NOW(),NOW()),
('user016','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser016','user016@test.com','13800000016',1,NOW(),NOW()),
('user017','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser017','user017@test.com','13800000017',1,NOW(),NOW()),
('user018','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser018','user018@test.com','13800000018',1,NOW(),NOW()),
('user019','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser019','user019@test.com','13800000019',1,NOW(),NOW()),
('user020','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy','TestUser020','user020@test.com','13800000020',1,NOW(),NOW());

SELECT COUNT(*) AS total_users FROM sys_user;
SELECT id,username,nickname,email FROM sys_user WHERE username LIKE 'user%' ORDER BY id LIMIT 10;
