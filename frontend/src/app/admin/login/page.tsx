'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import apiClient from '@/services/api';
import type { LoginRequest } from '@/types';

export default function AdminLoginPage() {
  const router = useRouter();
  const [formData, setFormData] = useState<LoginRequest>({
    email: '',
    password: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErrorMessage('');

    if (!formData.email || !formData.password) {
      setErrorMessage('Please enter both email and password');
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await apiClient.post('/admin/login', {
        email: formData.email.trim().toLowerCase(),
        password: formData.password,
      });

      if (response.success && response.data) {
        // Store token in API client
        apiClient.setToken(response.data.token);

        // Redirect to dashboard
        router.push('/admin/dashboard');
      }
    } catch (error: any) {
      if (error.code === 'INVALID_CREDENTIALS') {
        setErrorMessage('Invalid email or password');
      } else if (error.code === 'NETWORK_ERROR') {
        setErrorMessage('Unable to connect to server. Please check your connection.');
      } else {
        setErrorMessage(error.message || 'Login failed. Please try again.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-secondary via-secondary-light to-primary flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Logo/Title */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">
            Tau-Tau Run
          </h1>
          <p className="text-white/80 text-lg">
            Admin Dashboard
          </p>
        </div>

        {/* Login Card */}
        <div className="card">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">
            Sign In
          </h2>

          {/* Error Message */}
          {errorMessage && (
            <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg mb-6">
              <p className="font-medium">✗ {errorMessage}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Email Field */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                Email Address
              </label>
              <input
                type="email"
                id="email"
                value={formData.email}
                onChange={(e) => setFormData((prev) => ({ ...prev, email: e.target.value }))}
                className="input-field"
                placeholder="admin@tautaurun.com"
                disabled={isSubmitting}
                autoComplete="email"
                required
              />
            </div>

            {/* Password Field */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                Password
              </label>
              <input
                type="password"
                id="password"
                value={formData.password}
                onChange={(e) => setFormData((prev) => ({ ...prev, password: e.target.value }))}
                className="input-field"
                placeholder="••••••••"
                disabled={isSubmitting}
                autoComplete="current-password"
                required
              />
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isSubmitting}
              className="btn-secondary w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          {/* Default Credentials Hint (Development Only) */}
          {process.env.NODE_ENV === 'development' && (
            <div className="mt-6 p-4 bg-gray-50 rounded-lg border border-gray-200">
              <p className="text-xs text-gray-600 font-medium mb-2">
                Development Default Credentials:
              </p>
              <p className="text-xs text-gray-500 font-mono">
                Email: admin@tautaurun.com<br />
                Password: Admin123!
              </p>
            </div>
          )}
        </div>

        {/* Back to Home Link */}
        <div className="text-center mt-6">
          <a
            href="/"
            className="text-white hover:text-white/80 transition-colors text-sm"
          >
            ← Back to Home
          </a>
        </div>
      </div>
    </main>
  );
}
