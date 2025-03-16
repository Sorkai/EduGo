import axios from 'axios';

// 创建axios实例
const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器，添加token
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  role?: string;
  status?: string;
  createdAt?: string;
}

interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  message: string;
  token: string;
}

interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  firstName?: string;
  lastName?: string;
}

interface RegisterResponse {
  message: string;
  user: UserProfile;
}

interface UpdateUserRequest {
  email?: string;
  firstName?: string;
  lastName?: string;
}

interface ResetPasswordRequest {
  oldPassword: string;
  newPassword: string;
}

const userService = {
  // 用户登录
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/login', { username, password });
    return response.data;
  },
  
  // 用户注册
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await apiClient.post<RegisterResponse>('/register', data);
    return response.data;
  },
  
  // 获取用户信息
  getUserProfile: async (): Promise<UserProfile> => {
    const response = await apiClient.get<{ user: UserProfile }>('/user');
    return response.data.user;
  },
  
  // 更新用户信息
  updateUser: async (data: UpdateUserRequest): Promise<UserProfile> => {
    const response = await apiClient.put<{ message: string, user: UserProfile }>('/user', data);
    return response.data.user;
  },
  
  // 重置密码
  resetPassword: async (data: ResetPasswordRequest): Promise<{ message: string }> => {
    const response = await apiClient.put<{ message: string }>('/user/password', data);
    return response.data;
  },
  
  // 用户注销
  logout: async (): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>('/logout');
    // 清除本地存储的token
    localStorage.removeItem('token');
    sessionStorage.removeItem('token');
    return response.data;
  },
  
  // 刷新token
  refreshToken: async (): Promise<{ message: string, token: string }> => {
    const response = await apiClient.post<{ message: string, token: string }>('/refresh');
    return response.data;
  }
};

export default userService;
