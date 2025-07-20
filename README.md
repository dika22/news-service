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
git clone https://github.com/username/nama-proyek.git
cd nama-proyek

# set .env
cp -R .env.copy to .env

# generate swagger
swag ini

# Cara menjalankan http 
make http-serve

# Cara menjalankan worker
make start-worker
