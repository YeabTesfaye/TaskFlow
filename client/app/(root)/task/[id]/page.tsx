'use client';

import { useEffect, useState } from 'react';
import { notFound, useRouter } from 'next/navigation';
import { TaskEditor } from './TaskEditor';
import { getTask } from '@/lib/actions/task.action';
import { Task } from '@/types';
import { useAuth } from '@/hooks/use-auth';

interface EditTaskPageProps {
  params: { id: string };
}

const EditTaskPage = ({ params: { id } }: EditTaskPageProps) => {
  const [task, setTask] = useState<Task | null>(null);

  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push('/login');
    }
  }, [loading, isAuthenticated, router]);

  useEffect(() => {
    const fetchTask = async () => {
      const fetchedTask = await getTask(id);
      if (!fetchedTask) {
        notFound();
      } else {
        setTask(fetchedTask);
      }
    };

    fetchTask();
  }, [id]);

  if (!task) return <p>Loading...</p>;

  return <TaskEditor task={task} />;
};

export default EditTaskPage;
