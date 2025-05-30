import { useState, useEffect } from 'react';
import { statistics } from '@/api';

export interface TaskStats {
  total_tasks: number;
  completed_tasks: number;
  overdue_tasks: number;
  completion_rate: number;
  by_priority: Record<string, number>;
}

export const useTaskStats = (shouldFetch: boolean) => {
  const [stats, setStats] = useState<TaskStats | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!shouldFetch) return;

    const fetch = async () => {
      setLoading(true);
      try {
        const data = await statistics.getTaskStatistics();
        setStats(data);
      } catch (error) {
        console.error('Failed to fetch task stats', error);
      } finally {
        setLoading(false);
      }
    };

    fetch();
  }, [shouldFetch]);

  return { stats, loading };
};
