import { useEffect, useState, useCallback } from 'react';
import { getTasks } from '@/lib/actions/task.action';
import { Task } from '@/types';

export const useTasks = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [refetching, setRefetching] = useState(false);

  const fetchTasks = useCallback(async () => {
    setLoading(true);
    const fetched = await getTasks();
    setTasks(fetched);
    setLoading(false);
    setRefetching(false);
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  return {
    tasks,
    loading,
    refetching,
    refetch: fetchTasks,
  };
};
