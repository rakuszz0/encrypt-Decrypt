generateKey():
Membuat array/slice dari huruf-huruf alfabet.
Menggunakan algoritma Fisher-Yates shuffle untuk mengacak urutan huruf-huruf tersebut secara acak.
Menggabungkan kembali huruf-huruf yang sudah diacak menjadi sebuah string, yang kemudian menjadi kunci.
isValidKey():
Memvalidasi kunci yang diberikan oleh user.
Memastikan panjangnya 26 karakter.
Memastikan semua karakter adalah huruf dan tidak ada yang duplikat. Ini dilakukan dengan menggunakan map untuk melacak huruf yang sudah dilihat.
generateMapping() & generateReverseMapping():
generateMapping(): Membuat sebuah peta (map/dictionary) yang menghubungkan setiap huruf alfabet asli (a, b, c...) dengan huruf yang sesuai di kunci (x, p, j...). Ini adalah inti dari proses enkripsi.
generateReverseMapping(): Kebalikannya. Membuat peta dari huruf kunci kembali ke huruf alfabet asli. Ini diperlukan untuk dekripsi.
encryptText() & decryptText():
Kedua fungsi ini iterasi melalui setiap karakter dari teks input.
Jika karakter adalah huruf, ia akan mencarinya di peta (mapping atau reverseMapping) dan menggantinya.
Jika karakter bukan huruf (seperti spasi, angka, atau tanda baca), ia akan dibiarkan apa adanya.
5. Fungsi Handler API
Ini adalah fungsi yang menghubungkan logika inti dengan permintaan HTTP.

generateKeyHandler(c *gin.Context):
Memanggil generateKey().
Membungkus hasilnya dalam struct Response dan KeyResponse.
Mengirimkannya kembali ke client dalam format JSON dengan status 200 OK.
encryptHandler(c *gin.Context):
Bind JSON: c.ShouldBindJSON(&req) secara otomatis mem-parsing body JSON dari request dan memasukkannya ke dalam struct EncryptRequest. Jika formatnya salah, akan langsung mengembalikan error.
Validasi: Memanggil isValidKey() untuk memastikan kunci yang dikirim user valid.
Proses: Jika valid, memanggil encryptText() untuk melakukan enkripsi.
Respon: Mengirimkan hasil enkripsi dan peta pemetaannya (opsional, tapi bagus untuk demo) kembali ke client.
decryptHandler(c *gin.Context):
Strukturnya hampir sama dengan encryptHandler, tetapi menggunakan DecryptRequest dan memanggil decryptText().

