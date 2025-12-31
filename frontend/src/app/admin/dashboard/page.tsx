'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import apiClient from '@/services/api';

export default function AdminDashboardPage() {
  const router = useRouter();

  useEffect(() => {
    // Check if user is authenticated
    if (!apiClient.isAuthenticated()) {
      router.push('/admin/login');
    }
  }, [router]);

  const handleLogout = () => {
    apiClient.clearToken();
    router.push('/admin/login');
  };

  return (
    <main className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="container mx-auto px-4 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-800">
                Tau-Tau Run Admin
              </h1>
              <p className="text-sm text-gray-600">
                Event Management Dashboard
              </p>
            </div>
            <button
              onClick={handleLogout}
              className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="card max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">
            Welcome to Admin Dashboard
          </h2>
          <p className="text-gray-600 mb-6">
            You are successfully authenticated! Participant management features will be added in Phase 5.
          </p>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
            <h3 className="font-semibold text-blue-900 mb-2">
              ✅ Phase 4 Complete
            </h3>
            <p className="text-blue-800 text-sm">
              Admin authentication is working. You can log in and access protected routes.
            </p>
            <p className="text-blue-700 text-sm mt-3">
              <strong>Next:</strong> Phase 5 will add participant list and payment status management.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
