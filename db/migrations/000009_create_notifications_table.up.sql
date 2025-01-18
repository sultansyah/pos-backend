CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,  -- ID user yang menerima notifikasi (bisa admin atau kasir)
    title VARCHAR(255) NOT NULL,  -- Judul notifikasi
    message TEXT NOT NULL,  -- Pesan atau isi notifikasi
    status ENUM('unread', 'read') DEFAULT 'unread',  -- Status notifikasi
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
