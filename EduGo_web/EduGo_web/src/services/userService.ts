import axios from 'axios';
import config from '../config';

// 创建axios实例
const apiClient = axios.create({
  baseURL: config.apiBaseUrl,
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

export type UserRole = 'super_admin' | 'admin' | 'teacher' | 'student' | 'parent';

export const USER_ROLES = {
  SUPER_ADMIN: 'super_admin' as UserRole,
  ADMIN: 'admin' as UserRole,
  TEACHER: 'teacher' as UserRole,
  STUDENT: 'student' as UserRole,
  PARENT: 'parent' as UserRole
};

export const USER_ROLE_NAMES = {
  [USER_ROLES.SUPER_ADMIN]: '超级管理员',
  [USER_ROLES.ADMIN]: '管理员',
  [USER_ROLES.TEACHER]: '教师',
  [USER_ROLES.STUDENT]: '学生',
  [USER_ROLES.PARENT]: '家长'
};

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
  role?: UserRole;
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

// 用户管理相关接口
export interface UserRelation {
  id: number;
  userId: number;
  relatedUserId: number;
  relationType: string;
  status: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface AdminTeacherRelation extends UserRelation {
  department?: string;
  position?: string;
}

export interface TeacherStudentRelation extends UserRelation {
  courseId?: number;
  courseName?: string;
  semester?: string;
}

export interface StudentParentRelation extends UserRelation {
  relationship?: string;
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
  },
  
  // 用户管理相关
  
  // 获取所有用户（超级管理员权限）
  getAllUsers: async (): Promise<UserProfile[]> => {
    const response = await apiClient.get<{ users: UserProfile[] }>('/user-management/users');
    return response.data.users;
  },
  
  // 根据角色获取用户（管理员及以上权限）
  getUsersByRole: async (role: UserRole): Promise<UserProfile[]> => {
    const response = await apiClient.get<{ users: UserProfile[] }>(`/user-management/users/role/${role}`);
    return response.data.users;
  },
  
  // 根据ID获取用户（管理员及以上权限）
  getUserById: async (id: number): Promise<UserProfile> => {
    const response = await apiClient.get<{ user: UserProfile }>(`/user-management/users/${id}`);
    return response.data.user;
  },
  
  // 更新用户角色（超级管理员权限）
  updateUserRole: async (id: number, role: UserRole): Promise<{ message: string, user: { id: number, role: string } }> => {
    const response = await apiClient.put<{ message: string, user: { id: number, role: string } }>(`/super-admin/users/${id}/role`, { role });
    return response.data;
  },
  
  // 更新用户状态（管理员及以上权限）
  updateUserStatus: async (id: number, status: string): Promise<{ message: string, user: { id: number, status: string } }> => {
    const response = await apiClient.put<{ message: string, user: { id: number, status: string } }>(`/admin/users/${id}/status`, { status });
    return response.data;
  },
  
  // 用户关系管理
  
  // 创建管理员-教师关系（管理员及以上权限）
  createAdminTeacherRelation: async (teacherId: number, department?: string, position?: string): Promise<{ message: string, relation: AdminTeacherRelation }> => {
    const response = await apiClient.post<{ message: string, relation: AdminTeacherRelation }>('/admin/relations/teacher', { 
      teacher_id: teacherId,
      department,
      position
    });
    return response.data;
  },
  
  // 创建教师-学生关系（教师及以上权限）
  createTeacherStudentRelation: async (studentId: number, courseId?: number, courseName?: string, semester?: string): Promise<{ message: string, relation: TeacherStudentRelation }> => {
    const response = await apiClient.post<{ message: string, relation: TeacherStudentRelation }>('/teacher/relations/student', { 
      student_id: studentId,
      course_id: courseId,
      course_name: courseName,
      semester
    });
    return response.data;
  },
  
  // 创建学生-家长关系（学生及以上权限）
  createStudentParentRelation: async (parentId: number, relationship?: string): Promise<{ message: string, relation: StudentParentRelation }> => {
    const response = await apiClient.post<{ message: string, relation: StudentParentRelation }>('/student/relations/parent', { 
      parent_id: parentId,
      relationship
    });
    return response.data;
  },
  
  // 获取管理员管理的教师列表（管理员及以上权限）
  getTeachersByAdmin: async (): Promise<UserProfile[]> => {
    const response = await apiClient.get<{ teachers: UserProfile[] }>('/admin/relations/teachers');
    return response.data.teachers;
  },
  
  // 获取教师教授的学生列表（教师及以上权限）
  getStudentsByTeacher: async (): Promise<UserProfile[]> => {
    const response = await apiClient.get<{ students: UserProfile[] }>('/teacher/relations/students');
    return response.data.students;
  },
  
  // 获取学生的家长列表（学生及以上权限）
  getParentsByStudent: async (): Promise<UserProfile[]> => {
    const response = await apiClient.get<{ parents: UserProfile[] }>('/student/relations/parents');
    return response.data.parents;
  }
};

export default userService;
