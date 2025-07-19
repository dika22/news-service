CREATE TABLE articles (
    id SERIAL PRIMARY KEY,        -- Kolom 'id' sebagai primary key yang akan otomatis terisi dengan angka yang unik (auto-increment)
    author_id INT NOT NULL,       -- Kolom 'author_id' yang menunjukkan penulis. Misalnya ini adalah ID dari tabel 'authors'
    title VARCHAR(255) NOT NULL,  -- Kolom 'title' untuk judul postingan
    body TEXT NOT NULL,           -- Kolom 'body' untuk isi postingan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Kolom 'created_at' otomatis terisi dengan waktu sekarang
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Kolom 'updated_at' otomatis terisi dengan waktu sekarang
);
