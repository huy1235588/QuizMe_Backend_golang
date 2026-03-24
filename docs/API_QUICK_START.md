# QuizMe API - Quick Start Guide

## 🚀 Quick Start cho Frontend Developer

### 1. Base URL & Authorization

```javascript
const API_BASE_URL = 'http://localhost:8080/api';
const WS_URL = 'ws://localhost:8080/ws';

// Setup axios với token
import axios from 'axios';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Interceptor để thêm token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## 🔐 Authentication Flow

### 1. Đăng ký

```javascript
async function register(userData) {
  try {
    const response = await api.post('/auth/register', {
      username: userData.username,
      email: userData.email,
      password: userData.password,
      confirmPassword: userData.confirmPassword,
      fullName: userData.fullName
    });

    // Lưu tokens
    localStorage.setItem('accessToken', response.data.accessToken);
    localStorage.setItem('refreshToken', response.data.refreshToken);

    return response.data.user;
  } catch (error) {
    console.error('Register failed:', error.response.data);
    throw error;
  }
}
```

### 2. Đăng nhập

```javascript
async function login(usernameOrEmail, password) {
  try {
    const response = await api.post('/auth/login', {
      usernameOrEmail,
      password
    });

    localStorage.setItem('accessToken', response.data.accessToken);
    localStorage.setItem('refreshToken', response.data.refreshToken);

    return response.data.user;
  } catch (error) {
    console.error('Login failed:', error.response.data);
    throw error;
  }
}
```

### 3. Refresh Token

```javascript
async function refreshToken() {
  try {
    const refreshToken = localStorage.getItem('refreshToken');
    const response = await api.post('/auth/refresh-token', {
      refreshToken
    });

    localStorage.setItem('accessToken', response.data.accessToken);
    return response.data.accessToken;
  } catch (error) {
    // Nếu refresh token thất bại, logout user
    logout();
    throw error;
  }
}

// Setup interceptor để tự động refresh token khi expired
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshToken();
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);
```

### 4. Đăng xuất

```javascript
async function logout() {
  try {
    await api.post('/auth/logout');
  } finally {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    // Redirect to login page
  }
}
```

---

## 📚 Common API Usage Examples

### Quiz Management

```javascript
// Lấy danh sách quizzes
async function getQuizzes(page = 1, size = 10, filters = {}) {
  const params = new URLSearchParams({
    page,
    size,
    ...filters // search, categoryId, difficulty
  });

  const response = await api.get(`/quizzes/paged?${params}`);
  return response.data;
}

// Lấy chi tiết quiz
async function getQuizById(id) {
  const response = await api.get(`/quizzes/${id}`);
  return response.data;
}

// Tạo quiz mới
async function createQuiz(quizData) {
  const response = await api.post('/quizzes', {
    title: quizData.title,
    description: quizData.description,
    categoryIds: quizData.categoryIds,
    difficulty: quizData.difficulty, // EASY, MEDIUM, HARD
    isPublic: quizData.isPublic,
    questions: quizData.questions
  });
  return response.data;
}

// Lấy questions của quiz
async function getQuizQuestions(quizId) {
  const response = await api.get(`/questions/quiz/${quizId}`);
  return response.data;
}
```

### Room Management

```javascript
// Tạo room
async function createRoom(roomData) {
  const response = await api.post('/rooms', {
    name: roomData.name,
    quizId: roomData.quizId,
    maxPlayers: roomData.maxPlayers,
    password: roomData.password, // optional
    isPublic: roomData.isPublic
  });
  return response.data;
}

// Tham gia room bằng code
async function joinRoomByCode(code, guestName = null, password = null) {
  const response = await api.post('/rooms/join', {
    code,
    guestName,
    password
  });
  return response.data;
}

// Lấy danh sách rooms đang chờ
async function getWaitingRooms() {
  const response = await api.get('/rooms/waiting');
  return response.data;
}

// Bắt đầu game (host only)
async function startGame(roomId) {
  const response = await api.post(`/rooms/start/${roomId}`);
  return response.data;
}

// Rời room
async function leaveRoom(roomId) {
  const response = await api.delete(`/rooms/leave/${roomId}`);
  return response.data;
}
```

### User Management

```javascript
// Lấy profile hiện tại
async function getCurrentProfile() {
  const response = await api.get('/users/profile');
  return response.data;
}

// Upload avatar
async function uploadAvatar(file) {
  const formData = new FormData();
  formData.append('file', file);

  const response = await api.post('/users/avatar/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
  return response.data;
}

// Lấy top users
async function getTopUsers(limit = 10) {
  const response = await api.get(`/users/top?limit=${limit}`);
  return response.data;
}
```

### Category Management

```javascript
// Lấy tất cả categories
async function getAllCategories() {
  const response = await api.get('/categories');
  return response.data;
}

// Lấy active categories
async function getActiveCategories() {
  const response = await api.get('/categories/active');
  return response.data;
}
```

---

## 🔌 WebSocket Integration

### 1. Setup WebSocket Connection

```javascript
class GameWebSocket {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.messageHandlers = new Map();
  }

  connect(roomId, token, userId = null, guestName = null) {
    this.ws = new WebSocket(WS_URL);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;

      // Gửi JOIN message
      this.send('JOIN', {
        roomId,
        userId,
        guestName,
        token
      });
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.attemptReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  send(type, payload = null) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        type,
        payload,
        timestamp: new Date().toISOString()
      }));
    }
  }

  handleMessage(message) {
    const { type, payload, timestamp } = message;

    // Call registered handlers
    const handlers = this.messageHandlers.get(type) || [];
    handlers.forEach(handler => handler(payload, timestamp));
  }

  on(messageType, handler) {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, []);
    }
    this.messageHandlers.get(messageType).push(handler);
  }

  off(messageType, handler) {
    const handlers = this.messageHandlers.get(messageType) || [];
    const index = handlers.indexOf(handler);
    if (index > -1) {
      handlers.splice(index, 1);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      setTimeout(() => {
        console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        // Re-establish connection with same parameters
      }, 2000 * this.reconnectAttempts);
    }
  }
}
```

### 2. Usage trong Game Component

```javascript
// React example
import { useEffect, useState } from 'react';

function GameRoom({ roomId }) {
  const [ws] = useState(() => new GameWebSocket());
  const [gameState, setGameState] = useState({
    gameActive: false,
    currentQuestion: null,
    leaderboard: null
  });

  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    ws.connect(roomId, token);

    // Register message handlers
    ws.on('GAME_START', (payload) => {
      console.log('Game started!', payload);
      setGameState(prev => ({ ...prev, gameActive: true }));
    });

    ws.on('QUESTION', (payload) => {
      console.log('New question:', payload);
      setGameState(prev => ({
        ...prev,
        currentQuestion: payload
      }));
    });

    ws.on('LEADERBOARD', (payload) => {
      setGameState(prev => ({
        ...prev,
        leaderboard: payload
      }));
    });

    ws.on('QUESTION_RESULT', (payload) => {
      console.log('Question result:', payload);
      // Show results UI
    });

    ws.on('GAME_END', (payload) => {
      console.log('Game ended:', payload);
      setGameState(prev => ({ ...prev, gameActive: false }));
      // Show final results
    });

    ws.on('ERROR', (payload) => {
      console.error('WebSocket error:', payload);
      // Show error message
    });

    // Cleanup
    return () => {
      ws.disconnect();
    };
  }, [roomId]);

  const submitAnswer = (questionId, selectedOptions, answerTime) => {
    ws.send('ANSWER', {
      questionId,
      selectedOptions,
      answerTime
    });
  };

  const sendChatMessage = (message) => {
    ws.send('CHAT', {
      roomId,
      content: message
    });
  };

  return (
    <div>
      {gameState.currentQuestion && (
        <QuestionCard
          question={gameState.currentQuestion}
          onSubmit={submitAnswer}
        />
      )}

      {gameState.leaderboard && (
        <Leaderboard rankings={gameState.leaderboard.rankings} />
      )}
    </div>
  );
}
```

### 3. Answer Submission Example

```javascript
function QuestionCard({ question, onSubmit }) {
  const [selectedOptions, setSelectedOptions] = useState([]);
  const [startTime] = useState(Date.now());

  const handleSubmit = () => {
    const answerTime = (Date.now() - startTime) / 1000; // seconds
    onSubmit(question.questionId, selectedOptions, answerTime);
  };

  return (
    <div>
      <h3>{question.content}</h3>
      {question.imageUrl && <img src={question.imageUrl} alt="Question" />}

      <div>
        {question.options.map(option => (
          <button
            key={option.optionId}
            onClick={() => {
              if (question.questionType === 'CHECKBOX') {
                // Multiple selection
                setSelectedOptions(prev =>
                  prev.includes(option.optionId)
                    ? prev.filter(id => id !== option.optionId)
                    : [...prev, option.optionId]
                );
              } else {
                // Single selection
                setSelectedOptions([option.optionId]);
              }
            }}
            className={selectedOptions.includes(option.optionId) ? 'selected' : ''}
          >
            {option.content}
          </button>
        ))}
      </div>

      <button onClick={handleSubmit} disabled={selectedOptions.length === 0}>
        Submit Answer
      </button>
    </div>
  );
}
```

---

## 💡 Best Practices

### 1. Error Handling

```javascript
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response) {
      // Server responded with error status
      const { status, data } = error.response;

      switch (status) {
        case 400:
          // Validation error
          console.error('Validation error:', data.details);
          break;
        case 401:
          // Unauthorized - handled by refresh token interceptor
          break;
        case 403:
          // Forbidden
          console.error('You do not have permission');
          break;
        case 404:
          // Not found
          console.error('Resource not found');
          break;
        case 500:
          // Server error
          console.error('Server error:', data.error);
          break;
      }
    }

    return Promise.reject(error);
  }
);
```

### 2. Loading States

```javascript
function useApi(apiFunction) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [data, setData] = useState(null);

  const execute = async (...args) => {
    setLoading(true);
    setError(null);

    try {
      const result = await apiFunction(...args);
      setData(result);
      return result;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { loading, error, data, execute };
}

// Usage
const { loading, error, data, execute } = useApi(getQuizzes);

useEffect(() => {
  execute(1, 10, { difficulty: 'MEDIUM' });
}, []);
```

### 3. Caching với React Query

```javascript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

// Fetch quizzes
function useQuizzes(filters) {
  return useQuery({
    queryKey: ['quizzes', filters],
    queryFn: () => getQuizzes(filters.page, filters.size, filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Create quiz mutation
function useCreateQuiz() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createQuiz,
    onSuccess: () => {
      // Invalidate and refetch quizzes
      queryClient.invalidateQueries({ queryKey: ['quizzes'] });
    },
  });
}

// Usage
function QuizList() {
  const { data: quizzes, isLoading, error } = useQuizzes({ page: 1, size: 10 });
  const createMutation = useCreateQuiz();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {quizzes.data.map(quiz => (
        <QuizCard key={quiz.id} quiz={quiz} />
      ))}
    </div>
  );
}
```

### 4. WebSocket Heartbeat

```javascript
// Giữ connection alive với ping/pong
startHeartbeat() {
  this.heartbeatInterval = setInterval(() => {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.send('PING');
    }
  }, 30000); // 30 seconds
}

stopHeartbeat() {
  if (this.heartbeatInterval) {
    clearInterval(this.heartbeatInterval);
  }
}
```

---

## 🧪 Testing APIs

### Sử dụng Postman/Thunder Client

1. Import API collection (file `postman_collection.json` trong thư mục `api/`)
2. Setup environment variables:
   - `BASE_URL`: `http://localhost:8080`
   - `ACCESS_TOKEN`: Token sau khi login

### Sử dụng curl

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "usernameOrEmail": "admin",
    "password": "password123"
  }'

# Get quizzes với token
curl -X GET http://localhost:8080/api/quizzes \
  -H "Authorization: Bearer YOUR_TOKEN"

# Create room
curl -X POST http://localhost:8080/api/rooms \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Room",
    "quizId": 1,
    "maxPlayers": 10,
    "isPublic": true
  }'
```

---

## 📝 Common Issues & Solutions

### Issue: Token expired
**Solution:** Implement automatic token refresh using interceptors (xem phần Authentication Flow)

### Issue: WebSocket keeps disconnecting
**Solution:**
- Implement heartbeat (ping/pong)
- Add reconnection logic
- Check network connection

### Issue: CORS errors
**Solution:** Đảm bảo frontend origin được thêm vào `config.yaml`:
```yaml
cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:5173"
```

### Issue: 401 errors on protected endpoints
**Solution:**
- Check if token is being sent in Authorization header
- Verify token hasn't expired
- Ensure token format is `Bearer <token>`

---

## 📚 Additional Resources

- Full API Documentation: `API_DOCUMENTATION.md`
- Postman Collection: `api/postman_collection.json`
- WebSocket Events Reference: Section 10 trong API Documentation
- Example Frontend: (Coming soon)

---

**Happy Coding! 🚀**
