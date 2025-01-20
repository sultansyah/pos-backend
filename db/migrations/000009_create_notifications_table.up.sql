CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,  -- Judul notifikasi
    type VARCHAR(255) NOT NULL,  -- Type notifikasi seperti stock, transaction
    message TEXT NOT NULL,  -- Pesan atau isi notifikasi
    status ENUM('unread', 'read') DEFAULT 'unread',  -- Status notifikasi
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
