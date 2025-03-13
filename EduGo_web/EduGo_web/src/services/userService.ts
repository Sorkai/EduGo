export interface UserProfile {
  username: string;
  email: string;
  createdAt: string;
}

interface LoginResponse {
  token: string;
}

interface RegisterResponse {
  success: boolean;
}

const userService = {
  login: (email: string, password: string): Promise<LoginResponse> => {
    return Promise.resolve({ token: 'mock-token' });
  },
  register: (email: string, password: string): Promise<RegisterResponse> => {
    return Promise.resolve({ success: true });
  },
  getUserProfile: (): Promise<UserProfile> => {
    return Promise.resolve({
      username: 'MockUser',
      email: 'mock@example.com',
      createdAt: '2025-03-13'
    });
  }
};

export default userService;
