# News Service

Deskripsi singkat tentang proyek Golang ini.

## 🚀 Fitur

- Create Article
- Get Article

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

# how run swagger
http://localhost:8001/swagger/index.html
note : sesuaikan alamat url
```

### How run use docker 
```bash
# run use docker
docker compose up -d