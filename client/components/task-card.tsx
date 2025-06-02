'use client';

import { motion } from 'framer-motion';
import { CardHeaderSection } from './CardHeaderSection';
import { CardContentSection } from './CardContentSection';
import { CardFooterSection } from './CardFooterSection';
import { Task } from '@/types';
import { Card } from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';
import { useTags } from '@/hooks/useTags';
import { tasks } from '@/api';
import { useState } from 'react';

// Add to your imports
import { Users } from 'lucide-react';

interface TaskCardProps {
  task: Task;
  index: number;
  onDelete(): void;
  onStatusChange(status: Task['status']): void;
}

export function TaskCard({
  task,
  index,
  onDelete,
  onStatusChange,
}: TaskCardProps) {
  const { toast } = useToast();
  const { tagList } = useTags();
  const [localStatus, setLocalStatus] = useState<Task['status']>(task.status);

  const taskTags = tagList.filter((tag) => task.tags.includes(tag.id));

  const handleStatusChange = async (status: Task['status']) => {
    try {
      await tasks.updateStatus(task.id, { status });
      setLocalStatus(status);
      onStatusChange(status);
      toast({
        title: 'Task updated',
        description: `Task status changed to ${status}`,
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to update task status',
        variant: 'destructive',
        duration: 1000,
      });
    }
  };

  const handleDelete = async () => {
    try {
      await tasks.delete(task.id);
      toast({
        title: 'Task deleted',
        description: 'The task has been removed',
        duration: 1000,
      });
      onDelete();
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.message || 'Failed to delete task',
        variant: 'destructive',
        duration: 1000,
      });
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
    >
      <Card className="task-card">
        <CardHeaderSection
          task={task}
          onDelete={handleDelete}
          onStatusChange={handleStatusChange}
        />
        <CardContentSection task={task} />
        <CardFooterSection
          task={task}
          status={localStatus}
          tags={taskTags}
          onStatusChange={handleStatusChange}
        />
        <div className="flex items-center gap-2">
          {task.collaborators?.length > 0 && (
            <div className="flex items-center gap-1 text-muted-foreground">
              <Users size={16} />
              <span className="text-sm">{task.collaborators.length}</span>
            </div>
          )}
        </div>
      </Card>
    </motion.div>
  );
}

// Add to your TaskCard component's return JSX
