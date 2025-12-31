'use client';

import { useState, FormEvent } from 'react';
import apiClient from '@/services/api';
import type { RegisterRequest } from '@/types';

interface RegistrationFormProps {
  onSuccess?: () => void;
}

export default function RegistrationForm({ onSuccess }: RegistrationFormProps) {
  const [formData, setFormData] = useState<RegisterRequest>({
    name: '',
    email: '',
    phone: '',
    instagram_handle: '',
    address: '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name || formData.name.trim().length < 2) {
      newErrors.name = 'Name must be at least 2 characters';
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!formData.email || !emailRegex.test(formData.email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    if (!formData.phone || formData.phone.trim().length < 10) {
      newErrors.phone = 'Phone number must be at least 10 digits';
    }

    if (!formData.address || formData.address.trim().length < 10) {
      newErrors.address = 'Address must be at least 10 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage('');
    setErrorMessage('');
    setErrors({});

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await apiClient.post('/public/register', {
        name: formData.name.trim(),
        email: formData.email.trim().toLowerCase(),
        phone: formData.phone.trim(),
        instagram_handle: formData.instagram_handle?.trim() || undefined,
        address: formData.address.trim(),
      });

      if (response.success) {
        setSuccessMessage(
          response.message || 'Registration successful! Your payment status is pending.'
        );
        
        // Clear form
        setFormData({
          name: '',
          email: '',
          phone: '',
          instagram_handle: '',
          address: '',
        });

        if (onSuccess) {
          onSuccess();
        }
      }
    } catch (error: any) {
      if (error.code === 'DUPLICATE_EMAIL') {
        setErrorMessage('This email address is already registered.');
        setErrors({ email: 'Email already registered' });
      } else if (error.code === 'VALIDATION_ERROR' && error.details) {
        const validationErrors: Record<string, string> = {};
        error.details.forEach((detail: any) => {
          validationErrors[detail.field] = detail.message;
        });
        setErrors(validationErrors);
        setErrorMessage('Please check the form for errors.');
      } else {
        setErrorMessage(error.message || 'Registration failed. Please try again.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (field: keyof RegisterRequest, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field when user starts typing
    if (errors[field]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Success Message */}
      {successMessage && (
        <div className="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
          <p className="font-medium">✓ {successMessage}</p>
        </div>
      )}

      {/* Error Message */}
      {errorMessage && (
        <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
          <p className="font-medium">✗ {errorMessage}</p>
        </div>
      )}

      {/* Name Field */}
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
          Full Name <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="name"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          className={`input-field ${errors.name ? 'border-red-500' : ''}`}
          placeholder="John Doe"
          disabled={isSubmitting}
        />
        {errors.name && <p className="text-red-500 text-sm mt-1">{errors.name}</p>}
      </div>

      {/* Email Field */}
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
          Email Address <span className="text-red-500">*</span>
        </label>
        <input
          type="email"
          id="email"
          value={formData.email}
          onChange={(e) => handleChange('email', e.target.value)}
          className={`input-field ${errors.email ? 'border-red-500' : ''}`}
          placeholder="john.doe@example.com"
          disabled={isSubmitting}
        />
        {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email}</p>}
      </div>

      {/* Phone Field */}
      <div>
        <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-2">
          Phone Number <span className="text-red-500">*</span>
        </label>
        <input
          type="tel"
          id="phone"
          value={formData.phone}
          onChange={(e) => handleChange('phone', e.target.value)}
          className={`input-field ${errors.phone ? 'border-red-500' : ''}`}
          placeholder="+62 812 3456 7890"
          disabled={isSubmitting}
        />
        {errors.phone && <p className="text-red-500 text-sm mt-1">{errors.phone}</p>}
      </div>

      {/* Instagram Handle Field */}
      <div>
        <label htmlFor="instagram" className="block text-sm font-medium text-gray-700 mb-2">
          Instagram Handle <span className="text-gray-400">(Optional)</span>
        </label>
        <div className="relative">
          <span className="absolute left-3 top-2.5 text-gray-400">@</span>
          <input
            type="text"
            id="instagram"
            value={formData.instagram_handle}
            onChange={(e) => handleChange('instagram_handle', e.target.value)}
            className={`input-field pl-8 ${errors.instagram_handle ? 'border-red-500' : ''}`}
            placeholder="johndoe"
            disabled={isSubmitting}
          />
        </div>
        {errors.instagram_handle && (
          <p className="text-red-500 text-sm mt-1">{errors.instagram_handle}</p>
        )}
      </div>

      {/* Address Field */}
      <div>
        <label htmlFor="address" className="block text-sm font-medium text-gray-700 mb-2">
          Address <span className="text-red-500">*</span>
        </label>
        <textarea
          id="address"
          value={formData.address}
          onChange={(e) => handleChange('address', e.target.value)}
          className={`input-field min-h-[100px] ${errors.address ? 'border-red-500' : ''}`}
          placeholder="Jl. Example Street No. 123, Jakarta, Indonesia 12345"
          disabled={isSubmitting}
          rows={3}
        />
        {errors.address && <p className="text-red-500 text-sm mt-1">{errors.address}</p>}
      </div>

      {/* Submit Button */}
      <button
        type="submit"
        disabled={isSubmitting}
        className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? 'Registering...' : 'Register Now'}
      </button>

      <p className="text-sm text-gray-500 text-center">
        <span className="text-red-500">*</span> Required fields
      </p>
    </form>
  );
}
