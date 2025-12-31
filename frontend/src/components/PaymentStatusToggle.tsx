'use client';

import { useState } from 'react';

interface PaymentStatusToggleProps {
  participantId: string;
  currentStatus: 'PAID' | 'UNPAID';
  onUpdate: (id: string, newStatus: 'PAID' | 'UNPAID') => Promise<void>;
}

export default function PaymentStatusToggle({
  participantId,
  currentStatus,
  onUpdate,
}: PaymentStatusToggleProps) {
  const [isUpdating, setIsUpdating] = useState(false);
  const [localStatus, setLocalStatus] = useState(currentStatus);

  const handleToggle = async () => {
    const newStatus = localStatus === 'PAID' ? 'UNPAID' : 'PAID';
    
    // Optimistic update
    setLocalStatus(newStatus);
    setIsUpdating(true);

    try {
      await onUpdate(participantId, newStatus);
    } catch (error) {
      // Rollback on error
      setLocalStatus(currentStatus);
      console.error('Failed to update payment status:', error);
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <button
      onClick={handleToggle}
      disabled={isUpdating}
      className={`px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200 ${
        localStatus === 'PAID'
          ? 'bg-green-500 text-white hover:bg-green-600'
          : 'bg-gray-300 text-gray-700 hover:bg-gray-400'
      } ${isUpdating ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
    >
      {isUpdating ? (
        <span className="flex items-center gap-2">
          <svg
            className="animate-spin h-4 w-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          Updating...
        </span>
      ) : (
        localStatus
      )}
    </button>
  );
}
