CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    status ENUM('pending', 'completed', 'canceled') DEFAULT 'pending',
    code VARCHAR(50) NOT NULL UNIQUE,
    midtrans_status ENUM('pending', 'settlement', 'expire', 'cancel') DEFAULT 'pending',
    payment_status ENUM('unpaid', 'paid', 'failed') DEFAULT 'unpaid',
    payment_type ENUM('cash', 'qris', 'bank_transfer', 'credit_card') DEFAULT 'cash',
    payment_url VARCHAR(255),
    note TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE SET NULL
);