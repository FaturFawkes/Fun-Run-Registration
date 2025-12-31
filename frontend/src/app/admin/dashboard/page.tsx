'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import apiClient from '@/services/api';
import ParticipantList from '@/components/ParticipantList';
import type { Participant, ParticipantListResponse } from '@/types';

export default function AdminDashboardPage() {
  const router = useRouter();
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState({
    total: 0,
    paid: 0,
    unpaid: 0,
  });

  useEffect(() => {
    // Check if user is authenticated
    if (!apiClient.isAuthenticated()) {
      router.push('/admin/login');
      return;
    }

    // Fetch participants on mount
    fetchParticipants();
  }, [router]);

  const fetchParticipants = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await apiClient.get<ParticipantListResponse>('/admin/participants');
      
      if (response.success && response.data) {
        const participantsList = response.data.participants || [];
        setParticipants(participantsList);
        
        // Calculate stats
        const paid = participantsList.filter(p => p.payment_status === 'PAID').length;
        const unpaid = participantsList.filter(p => p.payment_status === 'UNPAID').length;
        
        setStats({
          total: response.data.total || 0,
          paid,
          unpaid,
        });
      }
    } catch (err: any) {
      console.error('Failed to fetch participants:', err);
      setError(err.message || 'Failed to load participants');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePaymentUpdate = async (id: string, newStatus: 'PAID' | 'UNPAID') => {
    try {
      const response = await apiClient.patch(`/admin/participants/${id}/payment`, {
        payment_status: newStatus,
      });

      if (response.success) {
        // Update local state
        setParticipants(prevParticipants =>
          prevParticipants.map(p =>
            p.id === id ? { ...p, payment_status: newStatus } : p
          )
        );

        // Update stats
        setStats(prevStats => {
          const adjustment = newStatus === 'PAID' ? 1 : -1;
          return {
            ...prevStats,
            paid: prevStats.paid + adjustment,
            unpaid: prevStats.unpaid - adjustment,
          };
        });
      }
    } catch (err: any) {
      console.error('Failed to update payment status:', err);
      throw err; // Re-throw to trigger rollback in PaymentStatusToggle
    }
  };

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
        {/* Statistics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-sm font-medium text-gray-500 uppercase">
              Total Participants
            </h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">{stats.total}</p>
          </div>
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-sm font-medium text-gray-500 uppercase">
              Paid
            </h3>
            <p className="text-3xl font-bold text-green-600 mt-2">{stats.paid}</p>
          </div>
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 className="text-sm font-medium text-gray-500 uppercase">
              Unpaid
            </h3>
            <p className="text-3xl font-bold text-orange-600 mt-2">{stats.unpaid}</p>
          </div>
        </div>

        {/* Participant List */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-xl font-bold text-gray-800">
              Participant List
            </h2>
            <p className="text-sm text-gray-600 mt-1">
              Manage participant payment status
            </p>
          </div>

          {error && (
            <div className="m-6 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-800 text-sm">{error}</p>
              <button
                onClick={fetchParticipants}
                className="mt-2 text-sm text-red-600 hover:text-red-800 font-medium"
              >
                Try Again
              </button>
            </div>
          )}

          <div className="p-6">
            <ParticipantList
              participants={participants}
              onPaymentUpdate={handlePaymentUpdate}
              isLoading={isLoading}
            />
          </div>
        </div>
      </div>
    </main>
  );
}
