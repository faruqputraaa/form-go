package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type Mahasiswa struct {
	gorm.Model
	NIM          string `form:"nim"`
	Nama         string `form:"nama"`
	TanggalLahir string `form:"tanggal_lahir"`
	NomorTelepon string `form:"nomor_telepon"`
	Email        string `form:"email"`
	Fakultas     string `form:"fakultas"`
	Prodi        string `form:"prodi"`
	Semester     int    `form:"semester"`
	JenisKelamin string `form:"jenis_kelamin"`
}

// Koneksi ke database MySQL menggunakan GORM
func connectDB() (*gorm.DB, error) {
	dsn := "root:my-secret-pw@tcp(127.0.0.1:3333)/?charset=utf8mb4&parseTime=True&loc=Local" // Ganti dengan user, password, dan nama database-mu
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

		// Buat database jika belum ada
		db.Exec("CREATE DATABASE IF NOT EXISTS mahasiswa")

		db.Exec("USE mahasiswa")

	// Auto migrate tabel mahasiswa
	db.AutoMigrate(&Mahasiswa{})

	return db, nil
}

func main() {
	// Inisialisasi Echo
	e := echo.New()

	// Koneksi ke database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Route untuk menampilkan form
	e.GET("/", showForm)

	// Route untuk menerima form submission
	e.POST("/submit", func(c echo.Context) error {
		return handleForm(c, db)
	})

	// Route untuk melayani file statis (CSS, JS, dll)
	e.Static("/src", "src")

	// Jalankan server
	e.Logger.Fatal(e.Start(":8080"))
}

// Fungsi untuk menampilkan form HTML
func showForm(c echo.Context) error {
	return c.File("templates/form.html")
}

// Fungsi untuk menangani form submission
func handleForm(c echo.Context, db *gorm.DB) error {
	// Ambil data dari form
	mahasiswa := new(Mahasiswa)
	if err := c.Bind(mahasiswa); err != nil {
		return c.String(http.StatusBadRequest, "Error parsing form")
	}

	// Validasi semester
	if mahasiswa.Semester <= 0 {
		return c.String(http.StatusBadRequest, "Semester harus lebih dari 0!")
	}

	// Simpan data ke database
	if err := db.Create(&mahasiswa).Error; err != nil {
		log.Println("Error menyimpan data:", err)
		return c.String(http.StatusInternalServerError, "Error menyimpan data ke database")
	}

	// Tampilkan data yang disubmit
	return c.String(http.StatusOK, fmt.Sprintf("Data Mahasiswa yang Dikirim: %+v\n", mahasiswa))
}
