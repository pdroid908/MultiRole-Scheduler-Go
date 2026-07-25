
package models

import(
	"time"
)

type Pengguna struct {
    Id           string    `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    PasswordHash string    `json:"password_hash" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Jadwal struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Tanggal     string    `json:"tanggal" db:"tanggal"`
	JamMulai    string    `json:"jam_mulai" db:"jam_mulai"`
	JamSelesai  string    `json:"jam_selesai" db:"jam_selesai"`
	Keterangan  *string   `json:"keterangan" db:"keterangan"`
	IsConfirmed bool      `json:"is_confirmed" db:"is_confirmed"`
}

type Jadwal_admin struct {
	
	UserId      string    `json:"user_id" db:"user_id"`
	Tanggal     string    `json:"tanggal" db:"tanggal"`
	JamMulai    string    `json:"jam_mulai" db:"jam_mulai"`
	JamSelesai  string    `json:"jam_selesai" db:"jam_selesai"`
	Keterangan  *string   `json:"keterangan" db:"keterangan"`
	IsConfirmed bool      `json:"is_confirmed" db:"is_confirmed"`
}
