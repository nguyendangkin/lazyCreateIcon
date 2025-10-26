Tuyệt vời, đây là file `README.md` cho chương trình của bạn:

```markdown
# Custom Desktop Entry Creator

Chương trình `custom-desktop-entry` giúp bạn tạo một shortcut (entry) trên menu ứng dụng của môi trường desktop Linux (như GNOME, KDE, XFCE, Cinnamon, v.v.) cho bất kỳ file thực thi nào, đồng thời tự động tải xuống và gán một icon tùy chỉnh từ một URL.

Điều này đặc biệt hữu ích cho các ứng dụng không được cài đặt thông qua trình quản lý gói truyền thống, hoặc các script cá nhân mà bạn muốn truy cập dễ dàng từ menu.

## Tính năng

*   **Tạo Desktop Entry:** Tự động tạo file `.desktop` theo chuẩn XDG.
*   **Tải Icon Tùy chỉnh:** Tải icon từ một URL được cung cấp và lưu trữ tại thư mục `~/.local/share/icons/custom-menu-icons/`.
*   **Hỗ trợ nhiều định dạng Icon:** Tự động nhận diện định dạng ảnh (PNG, JPG, SVG, v.v.) từ Content-Type hoặc phần mở rộng của URL.
*   **Tên an toàn:** Tự động làm sạch tên file thực thi để tạo tên entry an toàn.
*   **Cập nhật Menu:** Cố gắng gọi `update-desktop-database` để cập nhật menu ngay lập tức (nếu có trên hệ thống).

## Cách sử dụng

### Cài đặt (Build từ mã nguồn)

1.  Đảm bảo bạn đã cài đặt Go (phiên bản 1.16 trở lên).
2.  Clone hoặc tải xuống mã nguồn.
3.  Mở terminal trong thư mục chứa mã nguồn và chạy:
    ```bash
    go build -o custom-desktop-entry
    ```
    Thao tác này sẽ tạo ra một file thực thi có tên `custom-desktop-entry` trong cùng thư mục.
4.  (Tùy chọn) Di chuyển file thực thi vào một thư mục trong biến môi trường `PATH` của bạn để có thể chạy nó từ bất cứ đâu, ví dụ:
    ```bash
    sudo mv custom-desktop-entry /usr/local/bin/
    ```

### Chạy chương trình

Bạn cần cung cấp hai đối số:

1.  **Đường dẫn đầy đủ đến file thực thi:** Đây là ứng dụng hoặc script mà bạn muốn tạo shortcut cho.
2.  **URL của icon:** Một đường dẫn HTTP/HTTPS đến file hình ảnh (PNG, SVG, JPG, v.v.) sẽ được dùng làm icon.

**Cú pháp:**

```bash
custom-desktop-entry /duong/dan/day/du/den/ung_dung_cua_ban https://example.com/icon.png
```

**Ví dụ:**

Giả sử bạn có một script Python tại `/home/user/my-scripts/hello.py` và bạn muốn dùng icon từ `https://www.flaticon.com/svg/static/icons/svg/883/883907.svg`.

```bash
custom-desktop-entry /home/user/my-scripts/hello.py https://www.flaticon.com/svg/static/icons/svg/883/883907.svg
```

Hoặc nếu bạn đã tạo một file thực thi Go:

```bash
custom-desktop-entry /home/user/my-apps/my-go-app https://cdn-icons-png.flaticon.com/512/121/121543.png
```

## Giải thích chi tiết

Chương trình sẽ thực hiện các bước sau:

1.  **Kiểm tra đầu vào:** Xác nhận đường dẫn file thực thi và URL icon hợp lệ.
2.  **Tạo thư mục:**
    *   `~/.local/share/icons/custom-menu-icons/`: Nơi lưu trữ các icon đã tải về.
    *   `~/.local/share/applications/`: Thư mục chuẩn để lưu trữ các file `.desktop` do người dùng định nghĩa.
3.  **Tải Icon:**
    *   Tải file icon từ URL.
    *   Tự động suy luận phần mở rộng file (`.png`, `.svg`, `.jpg`, v.v.) từ Content-Type header hoặc từ URL.
    *   Lưu icon vào `~/.local/share/icons/custom-menu-icons/` với một tên file an toàn dựa trên tên của file thực thi. Nếu file icon đã tồn tại, nó sẽ thêm hậu tố số (`_1`, `_2`, v.v.) để tránh ghi đè.
4.  **Tạo Desktop Entry:**
    *   Tạo file `.desktop` tại `~/.local/share/applications/` với nội dung cơ bản:
        *   `Type=Application`
        *   `Name=` (lấy từ tên file thực thi đã được làm sạch)
        *   `Exec=` (đường dẫn đầy đủ đến file thực thi của bạn)
        *   `Icon=` (đường dẫn đầy đủ đến icon đã tải về)
        *   `Terminal=false` (ứng dụng sẽ không chạy trong terminal riêng)
        *   `Categories=Utility;` (bạn có thể chỉnh sửa thủ công sau nếu muốn)
        *   `StartupNotify=true`
    *   Gán quyền thực thi (0o755) cho file `.desktop` để đảm bảo hệ thống desktop có thể nhận diện nó.
5.  **Cập nhật Database:**
    *   Cố gắng chạy lệnh `update-desktop-database ~/.local/share/applications/`. Lệnh này giúp môi trường desktop phát hiện các entry mới mà không cần khởi động lại.
    *   Nếu lệnh này không tồn tại hoặc chạy không thành công, một cảnh báo sẽ hiển thị nhưng chương trình vẫn hoàn thành.

## Khắc phục sự cố

*   **Không thấy icon trong menu:**
    *   Thử đăng xuất và đăng nhập lại vào môi trường desktop của bạn.
    *   Khởi động lại panel hoặc dock của bạn.
    *   Đảm bảo `Exec` và `Icon` trong file `.desktop` trỏ đúng đến các file.
*   **Lỗi "executable not found":**
    *   Đảm bảo bạn đã cung cấp đường dẫn *đầy đủ* (absolute path) đến file thực thi, ví dụ: `/home/user/my-app` chứ không phải `~/my-app` hay `my-app`.
    *   Kiểm tra lại tên file và quyền thực thi của file đó.
*   **Lỗi tải icon:**
    *   Kiểm tra lại URL icon. Đảm bảo nó là một đường dẫn trực tiếp đến file ảnh.
    *   Kiểm tra kết nối mạng của bạn.
    *   Một số server có thể chặn yêu cầu từ user-agent mặc định; chương trình cố gắng gửi một user-agent đơn giản (`GoIconDownloader/1.0`).

## Góp ý và Đóng góp

Mọi góp ý hoặc đóng góp đều được hoan nghênh! Vui lòng mở một issue hoặc pull request trên kho lưu trữ của dự án.
```
