CREATE TABLE authors (
    id SERIAL PRIMARY KEY,        -- Kolom 'id' sebagai primary key yang akan otomatis terisi dengan angka yang unik (auto-increment)
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
