# QuizMe API Documentation

## 📌 Thông tin chung

- **Base URL**: `http://localhost:8080` (development)
- **API Version**: v1
- **Content-Type**: `application/json`
- **Authentication**: JWT Bearer Token

## 🔐 Authentication

API sử dụng JWT (JSON Web Token) để xác thực. Đối với các endpoint cần xác thực, thêm token vào header:

```
Authorization: Bearer <your_jwt_token>
```

### Token Types

- **Access Token**: Dùng cho các API requests, hết hạn sau 1 ngày (configurable)
- **Refresh Token**: Dùng để lấy access token mới, hết hạn sau 30 ngày (configurable)

---

## 📦 Response Format

Tất cả API responses đều được wrap trong một cấu trúc chung:

### Success Response

```json
{
    "status": "success",
    "message": "Operation successful message",
    "data": {
        // Actual response data here
    }
}
```

### Error Response

```json
{
    "status": "error",
    "message": "Error description",
    "data": null
}
```

### Paginated Response

```json
{
    "status": "success",
    "message": "Resources retrieved successfully",
    "data": {
        "content": [
            // Array of items
        ],
        "page": 1,
        "size": 10,
        "totalElements": 150,
        "totalPages": 15
    }
}
```

**Lưu ý:**

- Field `status` là **string** với giá trị `"success"` hoặc `"error"`
- Field `message` chứa thông báo mô tả kết quả
- Field `data` chứa dữ liệu thực tế hoặc `null` trong trường hợp lỗi
- Đối với danh sách được phân trang, `data` chứa một object với các field: `content`, `page`, `size`, `totalElements`, `totalPages`
- Khi không có dữ liệu, `data` có thể là `null` hoặc một object rỗng tùy thuộc vào endpoint
- Field `data` chứa dữ liệu thực tế (có thể là object, array, hoặc null)

Trong phần documention dưới đây, chỉ có phần `data` được mô tả chi tiết. Tất cả responses đều được wrap theo format trên.

---

## 📚 API Endpoints

### 1. Authentication APIs

#### 1.1. Đăng ký tài khoản

**POST** `/api/auth/register`

**Request Body:**

```json
{
    "username": "string (required, min=2, max=50)",
    "email": "string (required, valid email, max=100)",
    "password": "string (required, min=2, max=100)",
    "confirmPassword": "string (required, min=2, max=100)",
    "fullName": "string (required, max=100)"
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Registration successful",
    "data": {
        "accessToken": "string",
        "accessTokenExpiry": "2024-03-24T10:30:00Z",
        "refreshToken": "string",
        "refreshTokenExpiry": "2024-04-23T10:30:00Z",
        "user": {
            "id": 1,
            "username": "johndoe",
            "email": "john@example.com",
            "fullName": "John Doe",
            "profileImage": "https://...",
            "createdAt": "2024-03-24T10:30:00Z",
            "updatedAt": "2024-03-24T10:30:00Z",
            "lastLogin": "2024-03-24T10:30:00Z",
            "role": "USER",
            "isActive": true
        }
    }
}
```

#### 1.2. Đăng nhập

**POST** `/api/auth/login`

**Request Body:**

```json
{
    "usernameOrEmail": "string (required)",
    "password": "string (required)"
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "accessToken": "string",
        "accessTokenExpiry": "2024-03-24T10:30:00Z",
        "refreshToken": "string",
        "refreshTokenExpiry": "2024-04-23T10:30:00Z",
        "user": {
            "id": 1,
            "username": "johndoe",
            "email": "john@example.com",
            "fullName": "John Doe",
            "profileImage": "https://...",
            "createdAt": "2024-03-24T10:30:00Z",
            "updatedAt": "2024-03-24T10:30:00Z",
            "lastLogin": "2024-03-24T10:30:00Z",
            "role": "USER",
            "isActive": true
        }
    }
}
```

#### 1.3. Đăng xuất

**POST** `/api/auth/logout`

**Request Body:**

```json
{
    "refreshToken": "string (required)"
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Logout successful",
    "data": null
}
```

#### 1.4. Làm mới token

**POST** `/api/auth/refresh-token`

**Request Body:**

```json
{
    "refreshToken": "string (required)"
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Token refreshed successfully",
    "data": {
        "accessToken": "string",
        "accessTokenExpiry": "2024-03-24T10:30:00Z"
    }
}
```

---

### 2. User APIs

#### 2.1. Lấy thông tin user theo ID

**GET** `/api/users/:id`

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "User retrieved successfully",
    "data": {
        "id": 3,
        "username": "admin",
        "email": "admin@quizme.com",
        "fullName": "Administrator",
        "createdAt": "2026-03-22T16:30:38+07:00",
        "updatedAt": "2026-03-22T16:30:38+07:00",
        "role": "ADMIN",
        "isActive": true
    }
}
```

#### 2.2. Lấy top users theo điểm

**GET** `/api/users/top`

**Query Parameters:**

- `limit` (optional): Số lượng users (default: 10)

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Top users retrieved successfully",
    "data": [
        {
            "id": 1,
            "username": "johndoe",
            "email": "john@example.com",
            "fullName": "John Doe",
            "profileImage": "https://...",
            "createdAt": "2026-03-22T16:30:38+07:00",
            "updatedAt": "2026-03-22T16:30:38+07:00",
            "lastLogin": "2026-03-24T10:30:00+07:00",
            "role": "USER",
            "isActive": true
        }
    ]
}
```

#### 2.3. Đếm tổng số users

**GET** `/api/users/count`

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "User count retrieved successfully",
    "data": {
        "count": 1234
    }
}
```

#### 2.4. Lấy danh sách users (phân trang)

**GET** `/api/users/paged`

**Query Parameters:**

- `page` (optional): Trang hiện tại (default: 1)
- `size` (optional): Số items mỗi trang (default: 10)
- `search` (optional): Tìm kiếm theo username hoặc email

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Users retrieved successfully",
    "data": {
        "content": [
            {
                "id": 1,
                "username": "johndoe",
                "email": "john@example.com",
                "fullName": "John Doe",
                "profileImage": "https://...",
                "createdAt": "2026-03-22T16:30:38+07:00",
                "updatedAt": "2026-03-22T16:30:38+07:00",
                "lastLogin": "2026-03-24T10:30:00+07:00",
                "role": "USER",
                "isActive": true
            }
        ],
        "page": 1,
        "size": 10,
        "totalElements": 100,
        "totalPages": 10
    }
}
```

#### 2.5. Lấy profile user theo ID

**GET** `/api/users/profile/:id`

**Response:** `200 OK` - Same as get user by ID

#### 2.6. Lấy profile user hiện tại 🔒

**GET** `/api/users/profile`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK` - Same as get user by ID

#### 2.7. Upload avatar 🔒

**POST** `/api/users/avatar/upload`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**

```
file: <image file>
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Avatar uploaded successfully",
    "data": {
        "profileImage": "https://cloudinary.com/..."
    }
}
```

#### 2.8. Xóa avatar 🔒

**DELETE** `/api/users/avatar`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "Avatar deleted successfully",
    "data": null
}
```

#### 2.9. Tạo user mới 🔒👑

**POST** `/api/users/create`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Request Body:**

```json
{
    "username": "string",
    "email": "string",
    "password": "string",
    "fullName": "string",
    "role": "USER | ADMIN",
    "isActive": true
}
```

**Response:** `201 Created`

```json
{
    "status": "success",
    "message": "User created successfully",
    "data": {
        "id": 3,
        "username": "admin",
        "email": "admin@quizme.com",
        "fullName": "Administrator",
        "profileImage": null,
        "createdAt": "2026-03-22T16:30:38+07:00",
        "updatedAt": "2026-03-22T16:30:38+07:00",
        "lastLogin": null,
        "role": "ADMIN",
        "isActive": true
    }
}
```

#### 2.10. Cập nhật user 🔒👑

**PUT** `/api/users/:id`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Request Body:**

```json
{
    "username": "string (optional)",
    "email": "string (optional)",
    "fullName": "string (optional)",
    "password": "string (optional)",
    "profileImage": "string (optional)",
    "isActive": true
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "User updated successfully",
    "data": {
        "id": 3,
        "username": "admin",
        "email": "admin@quizme.com",
        "fullName": "Administrator",
        "profileImage": null,
        "createdAt": "2026-03-22T16:30:38+07:00",
        "updatedAt": "2026-03-22T16:30:38+07:00",
        "lastLogin": null,
        "role": "ADMIN",
        "isActive": true
    }
}
```

#### 2.11. Xóa user 🔒👑

**DELETE** `/api/users/:id`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "User deleted successfully",
    "data": null
}
```

#### 2.12. Khóa/Mở khóa user 🔒👑

**PUT** `/api/users/:id/lock`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Request Body:**

```json
{
    "isActive": false
}
```

**Response:** `200 OK`

```json
{
    "status": "success",
    "message": "User updated successfully",
    "data": {
        "id": 3,
        "username": "admin",
        "email": "admin@quizme.com",
        "fullName": "Administrator",
        "profileImage": null,
        "createdAt": "2026-03-22T16:30:38+07:00",
        "updatedAt": "2026-03-22T16:30:38+07:00",
        "lastLogin": null,
        "role": "ADMIN",
        "isActive": false
    }
}
```

---

### 3. Category APIs

#### 3.1. Lấy tất cả categories

**GET** `/api/categories`

**Response:** `200 OK`

```json
[
    {
        "id": 1,
        "name": "Science",
        "description": "Science related quizzes",
        "iconUrl": "https://...",
        "quizCount": 10,
        "totalPlayCount": 500,
        "isActive": true,
        "createdAt": "2024-03-24T10:30:00Z",
        "updatedAt": "2024-03-24T10:30:00Z"
    }
]
```

#### 3.2. Lấy category theo ID

**GET** `/api/categories/:id`

**Response:** `200 OK` - Single category object

#### 3.3. Lấy active categories

**GET** `/api/categories/active`

**Response:** `200 OK` - Array of active categories

#### 3.4. Tạo category mới 🔒👑

**POST** `/api/categories`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Request Body:**

```json
{
    "name": "string (required, max=100)",
    "description": "string (optional)",
    "iconUrl": "string (optional)"
}
```

**Response:** `201 Created`

#### 3.5. Cập nhật category 🔒👑

**PUT** `/api/categories/:id`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Request Body:**

```json
{
    "name": "string (optional)",
    "description": "string (optional)",
    "iconUrl": "string (optional)",
    "isActive": true
}
```

**Response:** `200 OK`

#### 3.6. Xóa category 🔒👑

**DELETE** `/api/categories/:id`

**Headers:** `Authorization: Bearer <token>` (Admin only)

**Response:** `200 OK`

---

### 4. Quiz APIs

#### 4.1. Lấy tất cả quizzes

**GET** `/api/quizzes`

**Response:** `200 OK`

```json
[
    {
        "id": 1,
        "title": "General Knowledge Quiz",
        "description": "Test your general knowledge",
        "quizThumbnails": "https://...",
        "categoryIds": [1, 2],
        "categoryNames": ["Science", "History"],
        "creatorId": 1,
        "creatorName": "John Doe",
        "creatorAvatar": "https://...",
        "difficulty": "MEDIUM",
        "isPublic": true,
        "playCount": 100,
        "questionCount": 10,
        "favoriteCount": 25,
        "createdAt": "2024-03-24T10:30:00Z",
        "updatedAt": "2024-03-24T10:30:00Z"
    }
]
```

#### 4.2. Lấy quiz theo ID

**GET** `/api/quizzes/:id`

**Response:** `200 OK` - Single quiz object

#### 4.3. Lấy public quizzes

**GET** `/api/quizzes/public`

**Response:** `200 OK` - Array of public quizzes

#### 4.4. Lấy quizzes theo độ khó

**GET** `/api/quizzes/difficulty/:difficulty`

**Path Parameters:**

- `difficulty`: `EASY` | `MEDIUM` | `HARD`

**Response:** `200 OK` - Array of quizzes

#### 4.5. Lấy quizzes (phân trang)

**GET** `/api/quizzes/paged`

**Query Parameters:**

- `page` (optional): Trang hiện tại (default: 1)
- `size` (optional): Số items mỗi trang (default: 10)
- `search` (optional): Tìm kiếm theo title
- `categoryId` (optional): Lọc theo category
- `difficulty` (optional): Lọc theo độ khó

**Response:** `200 OK`

```json
{
  "data": [...],
  "page": 1,
  "size": 10,
  "total": 100,
  "totalPages": 10
}
```

#### 4.6. Tạo quiz mới 🔒

**POST** `/api/quizzes`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
    "title": "string (required, max=100)",
    "description": "string (optional, max=1000)",
    "categoryIds": [1, 2],
    "difficulty": "EASY | MEDIUM | HARD (required)",
    "isPublic": true,
    "questions": [
        {
            "content": "string (required)",
            "imageUrl": "string (optional)",
            "videoUrl": "string (optional)",
            "audioUrl": "string (optional)",
            "funFact": "string (optional)",
            "explanation": "string (optional)",
            "timeLimit": 30,
            "points": 100,
            "orderNumber": 1,
            "type": "QUIZ",
            "options": [
                {
                    "content": "Option A",
                    "isCorrect": true
                },
                {
                    "content": "Option B",
                    "isCorrect": false
                }
            ]
        }
    ]
}
```

**Response:** `201 Created`

#### 4.7. Cập nhật quiz 🔒

**PUT** `/api/quizzes/:id`

**Headers:** `Authorization: Bearer <token>`

**Request Body:** Same as create quiz (all fields optional)

**Response:** `200 OK`

#### 4.8. Xóa quiz 🔒

**DELETE** `/api/quizzes/:id`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`

---

### 5. Question APIs

#### 5.1. Lấy question theo ID

**GET** `/api/questions/:id`

**Response:** `200 OK`

```json
{
    "id": 1,
    "quizId": 1,
    "content": "What is the capital of France?",
    "imageUrl": "https://...",
    "videoUrl": "https://...",
    "audioUrl": "https://...",
    "funFact": "Paris is known as the City of Light",
    "explanation": "Paris has been the capital since...",
    "timeLimit": 30,
    "points": 100,
    "orderNumber": 1,
    "type": "QUIZ",
    "options": [
        {
            "id": 1,
            "content": "Paris",
            "isCorrect": true
        },
        {
            "id": 2,
            "content": "London",
            "isCorrect": false
        }
    ],
    "createdAt": "2024-03-24T10:30:00Z",
    "updatedAt": "2024-03-24T10:30:00Z"
}
```

#### 5.2. Lấy questions theo quiz ID

**GET** `/api/questions/quiz/:quizId`

**Response:** `200 OK` - Array of questions

#### 5.3. Tạo question mới 🔒

**POST** `/api/questions`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
    "quizId": 1,
    "content": "string (required)",
    "imageUrl": "string (optional)",
    "videoUrl": "string (optional)",
    "audioUrl": "string (optional)",
    "funFact": "string (optional)",
    "explanation": "string (optional)",
    "timeLimit": 30,
    "points": 100,
    "orderNumber": 1,
    "type": "QUIZ | TRUE_FALSE | TYPE_ANSWER | QUIZ_AUDIO | QUIZ_VIDEO | CHECKBOX | POLL",
    "options": [
        {
            "content": "Option A",
            "isCorrect": true
        }
    ]
}
```

**Response:** `201 Created`

#### 5.4. Tạo nhiều questions cùng lúc 🔒

**POST** `/api/questions/batch`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "quizId": 1,
  "questions": [
    {
      "content": "Question 1",
      "timeLimit": 30,
      "points": 100,
      "orderNumber": 1,
      "type": "QUIZ",
      "options": [...]
    },
    {
      "content": "Question 2",
      "timeLimit": 30,
      "points": 100,
      "orderNumber": 2,
      "type": "QUIZ",
      "options": [...]
    }
  ]
}
```

**Response:** `201 Created`

#### 5.5. Cập nhật question 🔒

**PUT** `/api/questions/:id`

**Headers:** `Authorization: Bearer <token>`

**Request Body:** Same as create question (all fields optional)

**Response:** `200 OK`

#### 5.6. Xóa question 🔒

**DELETE** `/api/questions/:id`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`

---

### 6. Room APIs

#### 6.1. Lấy room theo code

**GET** `/api/rooms/:code`

**Response:** `200 OK`

```json
{
  "id": 1,
  "name": "Fun Quiz Room",
  "code": "ABC123",
  "quizId": 1,
  "hostId": 1,
  "quiz": {...},
  "host": {...},
  "hasPassword": false,
  "isPublic": true,
  "currentPlayerCount": 5,
  "maxPlayers": 10,
  "status": "WAITING",
  "startTime": null,
  "endTime": null,
  "createdAt": "2024-03-24T10:30:00Z",
  "participants": [
    {
      "id": 1,
      "user": {...},
      "score": 0,
      "isHost": true,
      "joinedAt": "2024-03-24T10:30:00Z",
      "leftAt": null,
      "isGuest": false,
      "guestName": null,
      "displayName": "John Doe"
    }
  ]
}
```

#### 6.2. Lấy waiting rooms

**GET** `/api/rooms/waiting`

**Response:** `200 OK` - Array of rooms with status WAITING

#### 6.3. Lấy available rooms

**GET** `/api/rooms/available`

**Query Parameters:**

- `limit` (optional): Số lượng rooms (default: 20)

**Response:** `200 OK` - Array of rooms

#### 6.4. Tạo room mới 🔒

**POST** `/api/rooms`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
    "name": "string (required, max=100)",
    "quizId": 1,
    "maxPlayers": 10,
    "password": "string (optional)",
    "isPublic": true
}
```

**Response:** `201 Created`

#### 6.5. Tham gia room bằng code 🔓

**POST** `/api/rooms/join`

**Headers:** `Authorization: Bearer <token>` (Optional - guests can join)

**Request Body:**

```json
{
    "code": "ABC123",
    "guestName": "Guest123 (optional)",
    "password": "string (optional)"
}
```

**Response:** `200 OK`

#### 6.6. Tham gia room bằng ID 🔓

**POST** `/api/rooms/join/:roomId`

**Headers:** `Authorization: Bearer <token>` (Optional)

**Request Body:**

```json
{
    "guestName": "Guest123 (optional)",
    "password": "string (optional)"
}
```

**Response:** `200 OK`

#### 6.7. Rời khỏi room 🔓

**DELETE** `/api/rooms/leave/:roomId`

**Headers:** `Authorization: Bearer <token>` (Optional)

**Response:** `200 OK`

#### 6.8. Đóng room 🔒

**PATCH** `/api/rooms/close/:roomId`

**Headers:** `Authorization: Bearer <token>` (Host only)

**Response:** `200 OK`

#### 6.9. Cập nhật room 🔒

**PATCH** `/api/rooms/:roomId`

**Headers:** `Authorization: Bearer <token>` (Host only)

**Request Body:**

```json
{
    "name": "string (optional)",
    "maxPlayers": 10,
    "password": "string (optional)",
    "isPublic": true
}
```

**Response:** `200 OK`

#### 6.10. Bắt đầu game 🔒

**POST** `/api/rooms/start/:roomId`

**Headers:** `Authorization: Bearer <token>` (Host only)

**Response:** `200 OK`

---

### 7. Chat APIs

#### 7.1. Lấy lịch sử chat

**GET** `/api/chat/room/:roomId`

**Query Parameters:**

- `limit` (optional): Số tin nhắn (default: 50)

**Response:** `200 OK`

```json
[
  {
    "id": 1,
    "roomId": 1,
    "user": {...},
    "isGuest": false,
    "guestName": null,
    "message": "Hello everyone!",
    "sentAt": "2024-03-24T10:30:00Z",
    "displayName": "John Doe"
  }
]
```

#### 7.2. Gửi tin nhắn 🔓

**POST** `/api/chat/send`

**Headers:** `Authorization: Bearer <token>` (Optional)

**Request Body:**

```json
{
    "roomId": 1,
    "content": "string (required, max=500)",
    "guestName": "Guest123 (optional)"
}
```

**Response:** `200 OK`

---

### 8. Game APIs

#### 8.1. Lấy trạng thái game 🔓

**GET** `/api/game/state/:roomId`

**Headers:** `Authorization: Bearer <token>` (Optional)

**Response:** `200 OK`

```json
{
    "gameActive": true,
    "currentQuestion": {
        "questionId": 1,
        "content": "What is the capital of France?",
        "imageUrl": "https://...",
        "questionType": "QUIZ",
        "timeLimit": 30,
        "points": 100,
        "options": [
            {
                "optionId": 1,
                "content": "Paris"
            },
            {
                "optionId": 2,
                "content": "London"
            }
        ]
    },
    "remainingTime": 25,
    "questionNumber": 1,
    "totalQuestions": 10,
    "leaderboard": {
        "rankings": [
            {
                "userId": 1,
                "username": "johndoe",
                "score": 500,
                "rank": 1,
                "isCorrect": true
            }
        ]
    }
}
```

#### 8.2. Khởi tạo game 🔒

**POST** `/api/game/init/:roomId`

**Headers:** `Authorization: Bearer <token>` (Host only)

**Response:** `200 OK`

#### 8.3. Bắt đầu game 🔒

**POST** `/api/game/start/:roomId`

**Headers:** `Authorization: Bearer <token>` (Host only)

**Response:** `200 OK`

---

### 9. Health Check

#### 9.1. Kiểm tra trạng thái server

**GET** `/health`

**Response:** `200 OK`

```json
{
    "status": "ok"
}
```

---

## 🔌 WebSocket API

### Connection

**Endpoint:** `ws://localhost:8080/ws`

### Message Format

All WebSocket messages follow this format:

```json
{
  "type": "MESSAGE_TYPE",
  "payload": {...},
  "timestamp": "2024-03-24T10:30:00Z"
}
```

### Message Types

#### Client → Server Messages

##### JOIN - Tham gia room

```json
{
    "type": "JOIN",
    "payload": {
        "roomId": 1,
        "userId": 1,
        "guestName": "Guest123",
        "token": "jwt_token"
    }
}
```

##### ANSWER - Gửi câu trả lời

```json
{
    "type": "ANSWER",
    "payload": {
        "questionId": 1,
        "selectedOptions": [1, 2],
        "answerTime": 15.5
    }
}
```

##### CHAT - Gửi tin nhắn

```json
{
    "type": "CHAT",
    "payload": {
        "roomId": 1,
        "content": "Hello!",
        "guestName": "Guest123"
    }
}
```

##### PING - Heartbeat

```json
{
    "type": "PING"
}
```

#### Server → Client Messages

##### GAME_START - Game bắt đầu

```json
{
    "type": "GAME_START",
    "payload": {
        "message": "Game is starting!"
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### QUESTION - Câu hỏi mới

```json
{
    "type": "QUESTION",
    "payload": {
        "questionId": 1,
        "content": "What is 2+2?",
        "imageUrl": null,
        "questionType": "QUIZ",
        "timeLimit": 30,
        "points": 100,
        "options": [
            { "optionId": 1, "content": "3" },
            { "optionId": 2, "content": "4" }
        ]
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### TIMER - Đếm ngược thời gian

```json
{
    "type": "TIMER",
    "payload": {
        "remainingTime": 25,
        "totalTime": 30
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### ANSWER_RESULT - Kết quả câu trả lời

```json
{
    "type": "ANSWER_RESULT",
    "payload": {
        "isCorrect": true,
        "score": 100
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### QUESTION_RESULT - Kết quả câu hỏi (sau khi tất cả trả lời)

```json
{
    "type": "QUESTION_RESULT",
    "payload": {
        "questionId": 1,
        "correctOptions": [2],
        "statistics": {
            "totalAnswers": 10,
            "correctCount": 7,
            "incorrectCount": 3,
            "avgTime": 12.5,
            "optionStats": [
                {
                    "optionId": 1,
                    "count": 3,
                    "percentage": 30.0
                },
                {
                    "optionId": 2,
                    "count": 7,
                    "percentage": 70.0
                }
            ]
        },
        "playerResults": [
            {
                "userId": 1,
                "username": "johndoe",
                "isCorrect": true,
                "score": 100,
                "answerTime": 10.5,
                "selectedOptions": [2]
            }
        ]
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### LEADERBOARD - Bảng xếp hạng

```json
{
    "type": "LEADERBOARD",
    "payload": {
        "rankings": [
            {
                "userId": 1,
                "username": "johndoe",
                "score": 500,
                "rank": 1,
                "isCorrect": true
            }
        ]
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### NEXT_QUESTION - Chuẩn bị câu hỏi tiếp theo

```json
{
    "type": "NEXT_QUESTION",
    "payload": {
        "questionNumber": 2
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### GAME_END - Game kết thúc

```json
{
    "type": "GAME_END",
    "payload": {
        "reason": "completed",
        "message": "Game finished!",
        "result": {
            "roomId": 1,
            "quizTitle": "General Knowledge",
            "totalQuestions": 10,
            "duration": 300,
            "finalRankings": [
                {
                    "userId": 1,
                    "username": "johndoe",
                    "totalScore": 850,
                    "correctAnswers": 9,
                    "avgAnswerTime": 12.5,
                    "rank": 1
                }
            ]
        }
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### PARTICIPANT - Thông tin người chơi thay đổi

```json
{
  "type": "PARTICIPANT",
  "payload": {
    "action": "joined | left",
    "participant": {...}
  },
  "timestamp": "2024-03-24T10:30:00Z"
}
```

##### CHAT - Tin nhắn chat

```json
{
    "type": "CHAT",
    "payload": {
        "userId": 1,
        "username": "johndoe",
        "message": "Hello!",
        "sentAt": "2024-03-24T10:30:00Z"
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### PLAYER_DISCONNECT - Người chơi mất kết nối

```json
{
    "type": "PLAYER_DISCONNECT",
    "payload": {
        "userId": 1,
        "username": "johndoe"
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### PLAYER_RECONNECT - Người chơi kết nối lại

```json
{
    "type": "PLAYER_RECONNECT",
    "payload": {
        "userId": 1,
        "username": "johndoe"
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### ERROR - Lỗi

```json
{
    "type": "ERROR",
    "payload": {
        "code": "INVALID_ROOM",
        "message": "Room not found"
    },
    "timestamp": "2024-03-24T10:30:00Z"
}
```

##### PONG - Heartbeat response

```json
{
    "type": "PONG",
    "timestamp": "2024-03-24T10:30:00Z"
}
```

---

## 📊 Data Models & Enums

### User Role

```
USER    - Người dùng thông thường
ADMIN   - Quản trị viên
```

### Difficulty

```
EASY    - Dễ
MEDIUM  - Trung bình
HARD    - Khó
```

### Room Status

```
WAITING     - Đang chờ
IN_PROGRESS - Đang chơi
COMPLETED   - Đã hoàn thành
CANCELLED   - Đã hủy
```

### Question Type

```
QUIZ        - Câu hỏi trắc nghiệm (chọn 1 đáp án)
TRUE_FALSE  - Đúng/Sai
TYPE_ANSWER - Nhập câu trả lời
QUIZ_AUDIO  - Trắc nghiệm có audio
QUIZ_VIDEO  - Trắc nghiệm có video
CHECKBOX    - Chọn nhiều đáp án
POLL        - Khảo sát (không có đáp án đúng/sai)
```

---

## ⚠️ Error Responses

All error responses follow this format:

```json
{
    "status": "error",
    "message": "Error description",
    "data": null
}
```

**Lưu ý:** Field `status` là string với giá trị `"error"`, KHÔNG phải boolean.

### Common HTTP Status Codes

- `200 OK` - Request thành công
- `201 Created` - Tạo resource thành công
- `400 Bad Request` - Request không hợp lệ (validation error)
- `401 Unauthorized` - Chưa đăng nhập hoặc token không hợp lệ
- `403 Forbidden` - Không có quyền truy cập
- `404 Not Found` - Resource không tồn tại
- `409 Conflict` - Conflict (ví dụ: username đã tồn tại)
- `500 Internal Server Error` - Lỗi server

### Example Error Response

**400 Bad Request:**

```json
{
    "status": "error",
    "message": "Invalid request body",
    "data": null
}
```

**401 Unauthorized:**

```json
{
    "status": "error",
    "message": "Invalid username/email or password",
    "data": null
}
```

**409 Conflict:**

```json
{
    "status": "error",
    "message": "Username already exists",
    "data": null
}
```

---

## 🔔 Notes

### Authentication

- 🔒 = Requires authentication (user logged in)
- 👑 = Requires admin role
- 🔓 = Optional authentication (guest allowed)

### Pagination

Most list endpoints support pagination with these query parameters:

- `page`: Page number (default: 1)
- `size`: Items per page (default: 10)
- `search`: Search keyword (optional)

### Timestamps

All timestamps are in ISO 8601 format (RFC3339): `2024-03-24T10:30:00Z`

### File Upload

File uploads use `multipart/form-data` encoding with:

- Max file size: 10MB (configurable)
- Supported formats: JPG, PNG, GIF, MP3, MP4 (depends on endpoint)

### Rate Limiting

API có thể có rate limiting để tránh abuse (tùy cấu hình).

---

## 📞 Support

Nếu có vấn đề hoặc câu hỏi, vui lòng tạo issue trên GitHub repository.

---

**Version:** 1.0.0
**Last Updated:** 2026-03-24
