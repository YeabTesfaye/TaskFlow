'use client';

import { useEffect, useState } from 'react';
import { notFound } from 'next/navigation';
import { TaskEditor } from './TaskEditor';
import { getTask } from '@/lib/actions/task.action';
import { Task } from '@/types/task';

interface EditTaskPageProps {
  params: { id: string };
}

const EditTaskPage = ({ params: { id } }: EditTaskPageProps) => {
  const [task, setTask] = useState<Task | null>(null);

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
