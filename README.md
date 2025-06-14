# Test

Tech Stack:
- `Docker` (for Postgres database), if testing is to be done without docker, then feel free to copy the migration scripts from `sql/migrations`
- `Gin` - HTTP Framework for Golang

## Running Instructions

```bash 
# fire up the DB
make docker-up
make docker-down

# runs the application
go run-local

# clear DB
make clean

# view the log
docker-compose logs --follow
tail -f log/app.log
```

Important file edits:
- `setting/setting.yaml` (contains server & database config)
- `docker-compose.yml` (contains docker database image config)

## Notes For The Simple System

### Regarding handling Atomic Operation on The Transactions

Untuk memastikan bahwa setiap transaksi terjadi secara atomik, salah satunya dapat dilakukan dengan menggunakan fitur transaction lock yang ada pada Database System yang digunakan, contohnya pada MySQL menggunakan `BEGIN TRANSCATION`, lalu kemudian setelah proses pengubahan nominal selesai untuk 1 operasi tersebut, dilakukan _commit_ . Selama _lock_ dari _transaction_ tersebut belum dilepaskan, maka operasi yang lain akan ter-pending.

Cara yang kedua adalah secara _programmatic_ dengan menerapkan sebuah __locking balance__ / __2-step operation__, dimana untuk setiap operasi memiliki **request** dan **confirm**. Request dan confirm tersebut terhubung melalui 1 referensi yang sama yaitu ID transaksi. Saldo user baru akan benar-benar berubah ketika langkah **confirm** sukses.

### Regarding potential race conditions

Dalam system yang di desain menggunakan Golang dan mekanisme transaction dari DBMS, sangat kecil kemungkinan untuk terjadi _race condition_. Kecuali, apabila pada implementasi _programmatic_ yang telah disebutkan diatas tidak dilakukan dengan benar, maka _race condition_ dapat terjadi. Kemungkinan kedua adalah penggunaan _goroutine_ yang tidak benar.

Selain itu, pool connection ke database juga harus di fine-tune untuk memastikan traffic yang masuk tetap berjalan lancar walaupun terjadi burst request.

Testing untuk skenario serupa dapat dilihat pada file: `internal/api/api_test.go`. Test dapat dijalankan dengan menggunakan:

```bash
go test -v -count 1 ./internal/api
```

Catatan:
- Setiap eksekusi test akan membuat 1 user unique
- Setiap eksekusi test akan menjalankan 100 goroutine, 100x transaksi, dengan nominal 1000.00

### Regarding Rollback Should a Failure Occurs In The Middle of a Transaction

Untuk mekanisme rollback yang dapat dilakukan __mid-transaction__, dapat menggunakan fitur transaction yang sama yang ada pada DBMS, dengan menggunakan `ROLLBACK`, semua perubahan yang terjadi saat proses transaction akan di revert ke saat sebelum transaction dan __lock__ akan dilepas.
