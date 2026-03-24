# API Documentation Files

Thư mục này chứa các tài liệu API và công cụ hỗ trợ frontend developer.

## 📁 Files

### 1. `postman_collection.json`
Postman/Thunder Client collection chứa tất cả API endpoints.

**Cách sử dụng:**

#### Với Postman:
1. Mở Postman
2. Click **Import** > **Choose Files**
3. Select `postman_collection.json`
4. Collection sẽ xuất hiện trong sidebar

#### Với Thunder Client (VS Code):
1. Mở VS Code
2. Install extension **Thunder Client**
3. Click **Collections** > **Import**
4. Select `postman_collection.json`

#### Cấu hình:
- Collection variables:
  - `BASE_URL`: `http://localhost:8080` (có thể thay đổi cho staging/production)
  - `ACCESS_TOKEN`: Tự động được set sau khi login/register
  - `REFRESH_TOKEN`: Tự động được set sau khi login/register

#### Flow sử dụng:
1. Run request **Login** hoặc **Register**
2. Tokens sẽ tự động được lưu
3. Tất cả requests có 🔒 sẽ tự động dùng token
4. Nếu token hết hạn, run **Refresh Token**

## 📚 Related Documentation

- **Main API Documentation**: `../API_DOCUMENTATION.md` - Tài liệu API đầy đủ, chi tiết
- **Quick Start Guide**: `../API_QUICK_START.md` - Hướng dẫn tích hợp nhanh cho frontend

## 🚀 Quick Test

### 1. Test Health Check
```bash
curl http://localhost:8080/health
```

Expected: `{"status":"ok"}`

### 2. Test Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirmPassword": "password123",
    "fullName": "Test User"
  }'
```

### 3. Test Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "usernameOrEmail": "testuser",
    "password": "password123"
  }'
```

Save the `accessToken` from response.

### 4. Test Protected Endpoint
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🔗 WebSocket Testing

### Using websocat (CLI tool)
```bash
# Install websocat
cargo install websocat
# or: brew install websocat

# Connect to WebSocket
websocat ws://localhost:8080/ws

# Send JOIN message
{"type":"JOIN","payload":{"roomId":1,"token":"YOUR_TOKEN"},"timestamp":"2024-03-24T10:30:00Z"}
```

### Using Browser Console
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected');
  ws.send(JSON.stringify({
    type: 'JOIN',
    payload: {
      roomId: 1,
      token: 'YOUR_TOKEN'
    }
  }));
};

ws.onmessage = (event) => {
  console.log('Message:', JSON.parse(event.data));
};
```

## 📊 API Endpoints Summary

### Public Endpoints (No Auth)
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh-token`
- `GET /api/users/:id`
- `GET /api/users/top`
- `GET /api/categories`
- `GET /api/quizzes`
- `GET /api/rooms/:code`
- `GET /health`

### Protected Endpoints (Auth Required) 🔒
- `GET /api/users/profile`
- `POST /api/users/avatar/upload`
- `POST /api/quizzes`
- `POST /api/rooms`
- `POST /api/rooms/start/:roomId`

### Admin Only Endpoints 👑
- `POST /api/users/create`
- `PUT /api/users/:id`
- `DELETE /api/users/:id`
- `POST /api/categories`
- `PUT /api/categories/:id`
- `DELETE /api/categories/:id`

### Optional Auth Endpoints 🔓
- `POST /api/rooms/join`
- `POST /api/chat/send`
- `GET /api/game/state/:roomId`

## 💡 Tips

1. **Token Expiry**: Access token expires after 1 day. Use refresh token to get new access token.

2. **Environment Variables**: Tạo nhiều environments (Dev, Staging, Production) trong Postman:
   - Dev: `http://localhost:8080`
   - Staging: `https://staging-api.quizme.com`
   - Production: `https://api.quizme.com`

3. **Pre-request Scripts**: Collection đã có script tự động set token sau login/register.

4. **Testing WebSocket**: Dùng [Postman WebSocket](https://learning.postman.com/docs/sending-requests/websocket/websocket/) hoặc browser dev tools.

5. **Error Handling**: Tất cả errors trả về format:
   ```json
   {
     "error": "Error message",
     "code": "ERROR_CODE",
     "details": {...}
   }
   ```

## 🐛 Troubleshooting

### Port 8080 is already in use
```bash
# Check what's using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or change port in config.yaml
```

### CORS errors
Add your frontend origin to `config.yaml`:
```yaml
cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:5173"
```

### WebSocket connection failed
- Check firewall settings
- Ensure backend is running
- Verify WebSocket URL: `ws://` (not `wss://` for local)

### 401 Unauthorized
- Token might be expired, use refresh token endpoint
- Verify token format: `Bearer <token>`
- Check if endpoint requires authentication

## 📞 Support

For issues or questions:
- Check `../API_DOCUMENTATION.md` for detailed API specs
- Check `../API_QUICK_START.md` for integration examples
- Create issue on GitHub repository

---

**Happy Testing! 🎯**
