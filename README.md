# News Service

Deskripsi singkat tentang proyek Golang ini.

## 🚀 Fitur

- Create Article [status : draft]
- Update Article [status : publish]
- Get Articles
- Get Article by ID

## Implementation
- Ratelimit -> handle by ip per request second
- Singleflight -> if many request and cache expired, hanya ada 1 request dari banyaknya request untuk melakukan query ke DB

## 🛠️ Teknologi

Service ini didevelop dengan:

- [Go](https://golang.org/) versi 1.23
- Modul Go (`go mod`)
- Database: PostgreSQL, ElasticSearch
- Cache: Redis
- Queue: RabbitMQ
- Tambahan: Docker

## 🧑‍💻 How Run Service

```bash
# clone repository
git clone https://github.com/dika22/news-service
cd nama-proyek

# set .env
cp -R .env.copy to .env
create name db
import table

# generate swagger
swag init

# Cara menjalankan http 
make http-serve

# Cara menjalankan worker
make start-worker

# how run unit test
make test

# how run swagger port sesuaikan dengan yang di .env
http://localhost:8001/swagger/index.html
note : sesuaikan alamat url
```

### How run use docker 
```bash
# run use docker
# Jalankan HTTP server
docker-compose up news-service #sesuaikan nama app jika perlu

# Jalankan worker background
docker-compose up worker

# Jalankan test
docker-compose run test
```

## Diagram Activity
# Create Article
![alt text](https://github.com/dika22/news-service/blob/main/create_article.png)

# Update Article
![alt text](https://github.com/dika22/news-service/blob/main/update_article.png)