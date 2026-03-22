# Admin User Seeding Script

Script để tạo tài khoản admin trong backend QuizMe. Script này hỗ trợ việc tạo mới hoặc cập nhật tài khoản admin với các tùy chọn linh hoạt.

## Cách Sử Dụng

### 1. Tạo Admin User với Thông Tin Mặc Định

```bash
cd QuizMe_Backend_golang
go run cmd/seed/main.go
```

**Thông tin mặc định:**
- Username: `admin`
- Email: `admin@quizme.com`
- Password: `admin123`
- Full Name: `Administrator`

### 2. Tạo Admin User với Thông Tin Tùy Chỉnh

```bash
go run cmd/seed/main.go \
  -username=yourusername \
  -email=your@email.com \
  -password=yourpassword \
  -fullname="Your Full Name"
```

**Ví dụ:**
```bash
go run cmd/seed/main.go \
  -username=john_admin \
  -email=john@quizme.com \
  -password=SecurePass123 \
  -fullname="John Administrator"
```

### 3. Cập Nhật Admin User Hiện Có

Nếu admin user đã tồn tại, sử dụng flag `-update`:

```bash
go run cmd/seed/main.go \
  -username=admin \
  -password=NewPassword123 \
  -update
```

## Các Tùy Chọn (Flags)

| Flag | Giá Trị Mặc Định | Mô Tả |
|------|------------------|-------|
| `-username` | `admin` | Tên đăng nhập của admin user |
| `-email` | `admin@quizme.com` | Email của admin user |
| `-password` | `admin123` | Mật khẩu (sẽ được mã hóa bằng bcrypt) |
| `-fullname` | `Administrator` | Tên đầy đủ của admin |
| `-update` | `false` | Cập nhật admin user nếu đã tồn tại |

## Yêu Cầu

1. **Database đã được khởi tạo:** Script sẽ tự động chạy migrations
2. **Environment Variables (Nếu có config file):**
   - `DB_HOST`: Database host (mặc định: `localhost`)
   - `DB_PORT`: Database port (mặc định: `5432`)
   - `DB_NAME`: Database name
   - `DB_USER`: Database user
   - `DB_PASSWORD`: Database password

3. **Config file** (tùy chọn): Đặt file `config.yaml` trong thư mục QuizMe_Backend_golang

Ví dụ `config.yaml`:
```yaml
database:
  driver: postgres
  host: localhost
  port: "5432"
  user: quizme_user
  password: quizme_password
  name: quizme_db
  sslmode: disable
```

## Ví Dụ Sử Dụng

### Ví dụ 1: Tạo admin với tất cả thông tin mặc định
```bash
go run cmd/seed/main.go
```
Output:
```
✓ Admin user created successfully!
  Username: admin
  Email: admin@quizme.com
  Full Name: Administrator
  Role: ADMIN
✓ User profile created successfully!
```

### Ví dụ 2: Tạo admin với thông tin tùy chỉnh
```bash
go run cmd/seed/main.go \
  -username=superadmin \
  -email=superadmin@company.com \
  -password=SuperSecure123! \
  -fullname="Super Administrator"
```

### Ví dụ 3: Cập nhật mật khẩu admin hiện có
```bash
go run cmd/seed/main.go \
  -username=admin \
  -password=NewSecurePassword123! \
  -update
```

## Tính Năng

✅ **Tạo mới admin user** với các thông tin tùy chỉnh
✅ **Mã hóa mật khẩu** sử dụng bcrypt (cost: 10)
✅ **Tự động tạo user profile** khi tạo user mới
✅ **Cập nhật admin user** nếu đã tồn tại (với flag `-update`)
✅ **Kiểm tra trùng lặp** username và email
✅ **Đảm bảo** user là admin và đang hoạt động (isActive = true)

## Kiểm Tra Kết Quả

Sau khi chạy script thành công, bạn có thể:

1. **Đăng nhập vào ứng dụng** với username và password đã tạo/cập nhật
2. **Kiểm tra database:**
   ```sql
   SELECT id, username, email, role, is_active FROM "user" WHERE role = 'ADMIN';
   ```

## Khắc Phục Sự Cố

### Lỗi: "Failed to connect to database"
- Kiểm tra database connection settings trong `config.yaml` hoặc environment variables
- Đảm bảo PostgreSQL server đang chạy
- Kiểm tra DB_* environment variables

### Lỗi: "Admin user already exists"
- Sử dụng flag `-update` để cập nhật user hiện có
- Hoặc sử dụng `-username` khác để tạo admin user mới

### Lỗi: "Email is already registered"
- Email này đã được sử dụng bởi user khác
- Sử dụng email khác hoặc cập nhật user với flag `-update`

## Bảo Mật

⚠️ **Lưu ý:**
- Không commit mật khẩu vào version control
- Thay đổi mật khẩu mặc định trước khi deploy lên production
- Sử dụng mật khẩu mạnh (ít nhất 12 ký tự) cho production
- Chỉ chạy script này từ một máy an toàn

## Liên Quan

- **Auth Service:** `internal/service/auth_service.go`
- **User Repository:** `internal/repository/user_repository.go`
- **User Domain:** `internal/domain/user.go`
- **Role Enum:** `internal/domain/enums/role.go`
