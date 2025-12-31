import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';
import type { APIResponse } from '@/types';

class APIClient {
  private client: AxiosInstance;
  private token: string | null = null;

  constructor() {
    this.client = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
      headers: {
        'Content-Type': 'application/json',
      },
      timeout: 10000, // 10 seconds
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response: AxiosResponse<APIResponse>) => {
        return response;
      },
      (error: AxiosError<APIResponse>) => {
        // Handle specific error cases
        if (error.response) {
          const { status, data } = error.response;

          if (status === 401 || status === 403) {
            // Token expired or unauthorized
            this.clearToken();
            
            // Redirect to login if not already there
            if (typeof window !== 'undefined' && !window.location.pathname.includes('/admin/login')) {
              window.location.href = '/admin/login';
            }
          }

          // Return structured error
          return Promise.reject({
            status,
            code: data?.error?.code || 'UNKNOWN_ERROR',
            message: data?.error?.message || 'An unexpected error occurred',
            details: data?.error?.details,
          });
        } else if (error.request) {
          // Network error
          return Promise.reject({
            status: 0,
            code: 'NETWORK_ERROR',
            message: 'Unable to connect to server. Please check your internet connection.',
          });
        } else {
          // Other errors
          return Promise.reject({
            status: 0,
            code: 'REQUEST_ERROR',
            message: error.message || 'An error occurred while making the request',
          });
        }
      }
    );

    // Load token from localStorage if available
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token');
    }
  }

  // Set authentication token
  setToken(token: string) {
    this.token = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  }

  // Clear authentication token
  clearToken() {
    this.token = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  }

  // Get current token
  getToken(): string | null {
    return this.token;
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.token !== null;
  }

  // Generic GET request
  async get<T = any>(url: string, params?: any): Promise<APIResponse<T>> {
    const response = await this.client.get<APIResponse<T>>(url, { params });
    return response.data;
  }

  // Generic POST request
  async post<T = any>(url: string, data?: any): Promise<APIResponse<T>> {
    const response = await this.client.post<APIResponse<T>>(url, data);
    return response.data;
  }

  // Generic PATCH request
  async patch<T = any>(url: string, data?: any): Promise<APIResponse<T>> {
    const response = await this.client.patch<APIResponse<T>>(url, data);
    return response.data;
  }

  // Generic DELETE request
  async delete<T = any>(url: string): Promise<APIResponse<T>> {
    const response = await this.client.delete<APIResponse<T>>(url);
    return response.data;
  }
}

// Export singleton instance
const apiClient = new APIClient();
export default apiClient;
