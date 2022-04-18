package Controllers

import (
	"database/sql"
	"time"
)

type NullInt64 sql.NullInt64

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Pengguna struct {
	Email    string `json:"email"`
	Nama     string `json:"nama"`
	Password string `json:"password"`
	Tipe     int    `json:"tipe"`
}

type PenggunaResponse struct {
	Data Pengguna `json:"data"`
}

type ArrPenggunaResponse struct {
	Data []Pengguna `json:"data"`
}

type Buku struct {
	Isbn       string    `json:"isbn"`
	Judul      string    `json:"judul"`
	Penulis    string    `json:"penulis"`
	Edisi      int       `json:"edisi"`
	TahunCetak time.Time `json:"tahun_cetak"`
	Harga      int       `json:"harga"`
	Genre      []Genre   `json:"genre"`
	PathFile   string    `json:"path_file"`
}

type BestSeller struct {
	Isbn             string  `json:"isbn`
	Judul            string  `json:"judul`
	Id_genre         int     `json:"id_genre`
	Genre            string  `json:"genre`
	Jumlah_penjualan int     `json:"jumlah`
	Buku             []Buku  `json:"buku"`
	Genres           []Genre `json:"genre"`
}

type BestSellerResponse struct {
	Data []BestSeller `json:"data"`
}

type Genre struct {
	IdGenre int    `json:"id_genre"`
	Genre   string `json:"genre"`
	Buku    []Buku `json:"buku"`
}

type GenreResponse struct {
	Data []Genre `json:"data"`
}

type BukuResponse struct {
	Data Buku `json:"data"`
}

type ArrBukuResponse struct {
	Data []Buku `json:"data"`
}

type UlasanPenilaian struct {
	Ulasan    string `json:"ulasan"`
	Penilaian int    `json:"penilaian"`
	Isbn      string `json:"isbn"`
	Email     string `json:"email"`
}

type UlasanPenilaianResponse struct {
	Data UlasanPenilaian `json:"data"`
}

type ArrUlasanPenilaianResponse struct {
	Data []UlasanPenilaian `json:"data"`
}

type DetailBuku struct {
	Isbn       string            `json:"isbn"`
	Judul      string            `json:"judul"`
	Penulis    string            `json:"penulis"`
	Edisi      int               `json:"edisi"`
	TahunCetak time.Time         `json:"tahun_cetak"`
	Harga      int               `json:"harga"`
	Genre      []Genre           `json:"genre"`
	Ulasan     []UlasanPenilaian `json:"ulasan"`
}

type DetailBukuResponse struct {
	Data DetailBuku `json:"data"`
}

type Kupon struct {
	Kode          int       `json:"kode"`
	BerlakuSampai time.Time `json:"berlaku_sampai"`
	Nominal       int       `json:"nominal"`
}

type KuponResponse struct {
	Data Kupon `json:"data"`
}

type Transaksi struct {
	IdTransaksi      int       `json:"id_transaksi"`
	Email            string    `json:"email"`
	Isbn             string    `json:"isbn"`
	TanggalTransaksi time.Time `json:"tanggal_transaksi"`
	JenisTransaksi   string    `json:"jenis_transaksi"`
	NominalTransaksi int       `json:"nominal_transaksi"`
	Buku             []Buku    `json:"buku"`
	Kupon            NullInt64 `json:"kupon"`
}

type TransaksiResponse struct {
	Data Transaksi `json:"data"`
}

type ArrTransaksiResponse struct {
	Data []Transaksi `json:"data"`
}

type Forum struct {
	IdForum      int       `json:"id_forum"`
	Email        string    `json:"email"`
	WaktuDikirim time.Time `json:"waktu_dikirim"`
	Pesan        string    `json:"pesan"`
}

type ForumResponse struct {
	Data Forum `json:"data"`
}

type ArrForumResponse struct {
	Data []Forum `json:"data"`
}
