import { tasks } from '@/api';
import { useState } from 'react';
import { Status, STATUS_VALUES } from '@/types';

interface StatusUpdaterProps {
  taskId: string;
  currentStatus: Status;
  onUpdated?: (newStatus: Status) => void;
}

export function StatusUpdater({
  taskId,
  currentStatus,
  onUpdated,
}: StatusUpdaterProps) {
  const [status, setStatus] = useState<Status>(currentStatus);
  const [loading, setLoading] = useState(false);

  const handleChange = async (newStatus: Status) => {
    setLoading(true);
    try {
      const updated = await tasks.updateStatus(taskId, { status: newStatus });
      setStatus(updated.status);
      onUpdated?.(updated.status);
    } catch (error) {
      console.error('Failed to update status:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <select
      value={status}
      onChange={(e) => handleChange(e.target.value as Status)}
      disabled={loading}
      className="border rounded px-2 py-1 text-sm"
    >
      {STATUS_VALUES.map((s) => (
        <option key={s} value={s}>
          {s}
        </option>
      ))}
    </select>
  );
}
