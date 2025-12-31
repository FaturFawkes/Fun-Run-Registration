// TypeScript type definitions for the application

export interface Participant {
  id: string;
  name: string;
  email: string;
  phone: string;
  instagram_handle: string | null;
  address: string;
  registration_status: 'PENDING' | 'CONFIRMED';
  payment_status: 'UNPAID' | 'PAID';
  created_at: string;
  updated_at: string;
}

export interface Admin {
  id: string;
  email: string;
}

export interface LoginResponse {
  token: string;
  admin: Admin;
  expires_at: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  phone: string;
  instagram_handle?: string;
  address: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  registration_status: string;
  payment_status: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface UpdatePaymentRequest {
  payment_status: 'PAID' | 'UNPAID';
}

export interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
}

export interface ValidationError {
  field: string;
  message: string;
}

export interface ParticipantListResponse {
  participants: Participant[];
  total: number;
  page: number;
  limit: number;
}
