
---

```markdown
# Linux App Shortcut Generator

Tự động tạo shortcut menu (file `.desktop`) và thêm icon tùy chỉnh cho ứng dụng bất kỳ trên **Linux Mint / Ubuntu / Debian**.
Chương trình viết bằng **Go**, chạy trực tiếp trong terminal.

---

## 🧩 Tính năng

- ✅ Tự động tải icon từ **link ảnh trên Internet** (PNG, JPG, SVG, v.v.)
- ✅ Tạo shortcut `.desktop` trong thư mục:
```

~/.local/share/applications/

```
- ✅ Lưu icon vào:
```

~/.local/share/icons/custom-menu-icons/

````
- ✅ Hỗ trợ tên file thực thi có khoảng trắng
- ✅ Không cần quyền root
- ✅ Gọi `update-desktop-database` để cập nhật menu tự động (nếu hệ thống có)

---

## ⚙️ Cài đặt & biên dịch

### 1. Cài Go (nếu chưa có)
```bash
sudo apt install golang
````

### 2. Clone hoặc tạo file `main.go`

Sao chép nội dung từ chương trình vào file `main.go`.

### 3. Biên dịch

```bash
go build -o program main.go
```

Sau khi build thành công, bạn sẽ có file thực thi `program`.

---

## 🚀 Cách sử dụng

Chạy lệnh:

```bash
./program /đường/dẫn/tới/file_thực_thi https://link.tới.ảnh.png
```

**Ví dụ:**

```bash
./program /home/ongchin/Tools/Chrome.AppImage https://upload.wikimedia.org/wikipedia/commons/a/a5/Google_Chrome_icon_(February_2022).svg
```

Khi chạy xong, chương trình sẽ:

* Tải icon về thư mục `~/.local/share/icons/custom-menu-icons/`
* Tạo file `.desktop` như sau:

  ```
  ~/.local/share/applications/Chrome.AppImage.desktop
  ```

---

## 📄 Nội dung file `.desktop` ví dụ

```ini
[Desktop Entry]
Type=Application
Name=Chrome_AppImage
Exec="/home/ongchin/Tools/Chrome.AppImage"
Icon="/home/ongchin/.local/share/icons/custom-menu-icons/Chrome_AppImage.svg"
Terminal=false
Categories=Utility;
StartupNotify=true
```

Sau khi tạo xong, ứng dụng sẽ xuất hiện trong menu “Start / Menu / Tất cả ứng dụng”.

---

## 🧰 Gỡ shortcut đã tạo

Xóa thủ công hai file:

```bash
rm ~/.local/share/applications/<tên>.desktop
rm ~/.local/share/icons/custom-menu-icons/<icon>.png
```

Sau đó, có thể chạy:

```bash
update-desktop-database ~/.local/share/applications
```

---

## ⚠️ Ghi chú

* Nếu icon chưa hiển thị, **đăng xuất / đăng nhập lại** hoặc **khởi động lại panel/menu**.
* Một số hệ thống yêu cầu `desktop-file-utils` để lệnh `update-desktop-database` hoạt động:

  ```bash
  sudo apt install desktop-file-utils
  ```
* Chương trình chỉ tạo shortcut **cho user hiện tại**, không cần `sudo`.

---

## 🧠 Bản quyền & Giấy phép

MIT License © 2025 — bạn có thể chỉnh sửa, phân phối hoặc tích hợp vào dự án của riêng mình.

```
