import apiClient from './api'

// Типы для примера
export interface User {
  id: number
  name: string
  email: string
}

export interface HealthCheckResponse {
  status: string
  service: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface RegisterResponse {
  id: string
  username: string
  email: string
  token: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
}

export interface ProfileResponse {
  user_id: string
  username: string
  email: string
}

// Пример API endpoints
export const apiEndpoints = {
  // Health check
  healthCheck: () => apiClient.get<HealthCheckResponse>('/health'),
  
  // Authentication
  register: (data: RegisterRequest) => apiClient.post<RegisterResponse>('/v1/users/register', data),
  login: (data: LoginRequest) => apiClient.post<LoginResponse>('/v1/users/login', data),
  
  // Users
  getProfile: () => apiClient.get<ProfileResponse>('/v1/users/profile'),
  getUsers: () => apiClient.get<User[]>('/users'),
  getUserById: (id: number) => apiClient.get<User>(`/users/${id}`),
  createUser: (user: Omit<User, 'id'>) => apiClient.post<User>('/users', user),
  updateUser: (id: number, user: Partial<User>) => apiClient.put<User>(`/users/${id}`, user),
  deleteUser: (id: number) => apiClient.delete(`/users/${id}`),
}
