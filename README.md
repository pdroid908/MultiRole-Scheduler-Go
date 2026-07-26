# 📖 Tentang Proyek

MultiRole Scheduler Go adalah aplikasi penjadwalan dan booking online yang dikembangkan menggunakan Go (Golang) untuk membantu penyedia layanan mengelola jadwal reservasi secara lebih efisien.

Aplikasi ini memungkinkan administrator mengatur jadwal dan data member, member mengelola layanan yang dimiliki, serta pelanggan melakukan booking secara online melalui tautan publik tanpa perlu memiliki akun administrator.

Setiap proses booking divalidasi oleh backend sehingga jadwal yang sudah dipesan tidak dapat dipilih kembali oleh pelanggan lain. Dengan mekanisme tersebut, sistem mampu mengurangi risiko bentrok jadwal (double booking) yang sering terjadi pada proses reservasi manual melalui telepon atau aplikasi chat..

# 💡 Latar Belakang

Banyak usaha kecil maupun penyedia jasa masih menerima reservasi melalui WhatsApp, telepon, atau pencatatan manual.

Seiring meningkatnya jumlah pelanggan, proses tersebut sering menimbulkan berbagai kendala, seperti:

- Jadwal yang sama dipesan oleh lebih dari satu pelanggan.
- Sulit mengetahui jadwal yang masih tersedia.
- Riwayat booking tidak terdokumentasi dengan baik.
- Admin harus memeriksa jadwal secara manual sebelum menerima pelanggan baru.
- Kesalahan pencatatan menyebabkan bentrok jadwal dan menurunkan kualitas pelayanan.

Permasalahan tersebut menjadi alasan dikembangkannya aplikasi ini.

# ✅ Solusi

MultiRole Scheduler Go menyediakan platform terpusat untuk mengelola proses reservasi.

Administrator dapat mengatur seluruh data sistem melalui dashboard, sedangkan member dapat melihat jadwal dan booking yang diterima.

Pelanggan tidak perlu membuat akun administrator. Cukup membuka tautan booking yang dibagikan, pelanggan dapat melihat jadwal yang tersedia dan melakukan reservasi secara online.

Seluruh proses booking divalidasi di sisi backend sehingga hanya jadwal yang masih tersedia yang dapat dipesan.

# 👥 Hak Akses Pengguna

## 👨‍💼 Administrator

Administrator bertanggung jawab mengelola seluruh sistem.

Fitur yang tersedia:

- Login ke dashboard
- Registrasi akun member
- Mengelola data member
- Mengelola jadwal
- Melihat seluruh data booking
- Membagikan tautan booking kepada pelanggan

---

## 👤 Member

Member melihat jadwal tanggal, waktu, dan keperluan, 
member hanya melihat jadwal yang sudah di validasi oleh admin dan resmi di terima untuk acara
---

## 🌐 Pelanggan

Pelanggan tidak perlu memiliki akses ke dashboard administrator.

Cukup dengan membuka **link booking** yang diberikan Admin, pelanggan dapat:

- Melihat jadwal yang tersedia
- Memilih jadwal
- Mengisi data pemesanan
- Mengirim permintaan booking

Setelah booking berhasil, jadwal tersebut otomatis masuk ke admin untuk di cek dan di setujui jika sesuai.

# 🔒 Pencegahan Bentrok Jadwal

Salah satu fitur utama dari aplikasi ini adalah mencegah terjadinya **double booking**.

Alur validasi yang dilakukan sistem:

1. Pelanggan memilih jadwal yang tersedia.
2. Permintaan booking dikirim ke server.
3. Backend memeriksa apakah jadwal masih tersedia.
4. Jika tersedia, data booking akan disimpan.
5. Jika jadwal sudah digunakan pelanggan lain, sistem akan menolak permintaan booking.

Dengan proses ini, satu jadwal hanya dapat dimiliki oleh satu pelanggan sehingga bentrok jadwal dapat dihindari.

# 🎯 Tujuan Pengembangan

Project ini dibuat sebagai implementasi backend menggunakan Go (Golang) dengan fokus pada:

- RESTful API
- JWT Authentication
- Role-Based Access Control (RBAC)
- Manajemen Jadwal
- Booking Online
- Validasi Double Booking
- Integrasi PostgreSQL
- Pengembangan aplikasi yang mendekati kebutuhan dunia kerja

# 🚀 Fitur Utama

- Autentikasi Login & Registrasi
- Role-Based Access Control (RBAC)
- Dashboard Administrator
- Dashboard Member
- CRUD Jadwal
- CRUD Data Member
- Booking Online
- Link Booking Publik
- Validasi Bentrok Jadwal
- Riwayat Booking
- REST API menggunakan Go (Gin)
- Database PostgreSQL