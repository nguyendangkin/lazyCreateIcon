package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: program /full/path/to/executable https://example.com/icon.png")
		os.Exit(2)
	}

	execPath, iconURL := os.Args[1], os.Args[2]

	// Kiểm tra file thực thi tồn tại
	absExec, err := filepath.Abs(execPath)
	if err != nil {
		exitErr(err)
	}
	if stat, err := os.Stat(absExec); err != nil || stat.IsDir() {
		exitErr(fmt.Errorf("executable not found or is a directory: %s", absExec))
	}

	home := os.Getenv("HOME")
	if home == "" {
		exitErr(errors.New("HOME environment variable not set"))
	}

	iconsDir := filepath.Join(home, ".local", "share", "icons", "custom-menu-icons")
	appsDir := filepath.Join(home, ".local", "share", "applications")

	// Tạo thư mục nếu chưa có
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		exitErr(err)
	}
	if err := os.MkdirAll(appsDir, 0o755); err != nil {
		exitErr(err)
	}

	// Tên an toàn cho entry: lấy từ tên file thực thi
	baseName := filepath.Base(absExec)
	safeName := sanitizeName(strings.TrimSuffix(baseName, filepath.Ext(baseName)))

	// Tải hình
	fmt.Println("Downloading icon:", iconURL)
	iconPath, err := downloadIcon(iconURL, iconsDir, safeName)
	if err != nil {
		exitErr(err)
	}
	fmt.Println("Saved icon to:", iconPath)

	// Tạo file .desktop
	desktopPath := filepath.Join(appsDir, safeName+".desktop")
	if err := writeDesktopFile(desktopPath, safeName, absExec, iconPath); err != nil {
		exitErr(err)
	}
	// Quyền để file .desktop có thể được nhận diện (nhiều distro chấp nhận 0644; một số yêu cầu executable)
	if err := os.Chmod(desktopPath, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "warning: chmod .desktop failed:", err)
	}

	fmt.Println("Created desktop entry:", desktopPath)

	// Cố gắng cập nhật database menu (nếu có)
	if err := tryUpdateDesktopDatabase(appsDir); err != nil {
		fmt.Fprintln(os.Stderr, "warning: update-desktop-database:", err)
	}

	fmt.Println("Done. Nếu bạn không thấy icon trong menu, thử đăng xuất/đăng nhập lại hoặc khởi động lại panel.")
}

// sanitizeName: giữ lại ký tự an toàn cho tên file/entry
func sanitizeName(s string) string {
	// chỉ giữ letters, numbers, '-', '_' và khoảng trắng (thay bằng _)
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_ ]+`)
	clean := re.ReplaceAllString(s, "")
	clean = strings.TrimSpace(clean)
	if clean == "" {
		return "app"
	}
	// thay space bằng dấu gạch dưới
	return strings.ReplaceAll(clean, " ", "_")
}

// downloadIcon tải ảnh từ URL, tự động chọn extension. Trả về đường dẫn file đã lưu.
func downloadIcon(url, destDir, baseName string) (string, error) {
	// timeout ngắn
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	// một số server cần header
	req.Header.Set("User-Agent", "GoIconDownloader/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to download icon: status %s", resp.Status)
	}

	// Thử lấy extension từ URL
	ext := filepath.Ext(resp.Request.URL.Path)
	ext = strings.ToLower(ext)
	if ext == "" || ext == "." {
		// Thử từ Content-Type
		ct := resp.Header.Get("Content-Type")
		if ct != "" {
			if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
				ext = exts[0]
			}
		}
	}
	// Fallback
	if ext == "" || ext == "." {
		ext = ".png"
	}

	// Nếu ext là svg+xml, mime.ExtensionsByType trả .svg ok.
	// Compose filename
	filename := baseName + ext
	full := filepath.Join(destDir, filename)

	// Nếu file tồn tại, thêm suffix để tránh ghi đè
	if _, err := os.Stat(full); err == nil {
		for i := 1; ; i++ {
			full = filepath.Join(destDir, fmt.Sprintf("%s_%d%s", baseName, i, ext))
			if _, err := os.Stat(full); os.IsNotExist(err) {
				break
			}
		}
	}

	out, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return full, nil
}

// writeDesktopFile tạo nội dung .desktop cơ bản
func writeDesktopFile(path, name, execPath, iconPath string) error {
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;
StartupNotify=true
`, name, escapeExec(execPath), escapeExec(iconPath))

	return os.WriteFile(path, []byte(content), 0o644)
}

// escapeExec: nếu đường dẫn có khoảng trắng thì bọc trong quotes (desktop spec cho phép)
// nhưng giữ Exec đơn giản (không adds args)
func escapeExec(s string) string {
	if strings.ContainsAny(s, " \t") {
		return `"` + s + `"`
	}
	return s
}

// tryUpdateDesktopDatabase: cố gắng chạy update-desktop-database nếu có
func tryUpdateDesktopDatabase(appsDir string) error {
	_, err := exec.LookPath("update-desktop-database")
	if err != nil {
		// không có lệnh này: không phải lỗi nghiêm trọng
		return err
	}
	cmd := exec.Command("update-desktop-database", appsDir)
	// set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		return fmt.Errorf("timeout running update-desktop-database")
	}
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
