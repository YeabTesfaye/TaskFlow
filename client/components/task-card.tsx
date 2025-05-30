'use client';

import { Task } from '@/types';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { PriorityBadge } from '@/components/ui/priority-badge';
import { StatusBadge } from '@/components/ui/status-badge';
import { TagBadge } from '@/components/ui/tag-badge';
import { formatDistanceToNow } from 'date-fns';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  MoreHorizontal,
  Pencil,
  Trash2,
  CheckCircle,
  RotateCw,
  Clock,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { useToast } from '@/hooks/use-toast';
import { motion } from 'framer-motion';
import { tasks } from '@/api';
import { useTags } from '@/hooks/useTags';
import { StatusUpdater } from './status-updater';
import { useState } from 'react';

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
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to update task status',
        variant: 'destructive',
      });
    }
  };

  const handleDelete = async () => {
    try {
      await tasks.delete(task.id);
      toast({
        title: 'Task deleted',
        description: 'The task has been removed',
      });
      onDelete();
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.message || 'Failed to delete task',
        variant: 'destructive',
      });
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
    >
      <Card className="h-full overflow-hidden rounded-2xl border border-border bg-background transition-all duration-200 hover:shadow-lg hover:ring-1 hover:ring-ring hover:scale-[1.01] active:scale-[0.99]">
        <CardHeader className="flex flex-row items-start justify-between gap-2 p-4 pb-2">
          <div className="space-y-1.5">
            <Link
              href={`/task/${task.id}`}
              className="line-clamp-2 text-base font-semibold leading-tight text-foreground hover:underline"
            >
              {task.title}
            </Link>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreHorizontal className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[180px]">
              <DropdownMenuItem asChild>
                <Link
                  href={`/task/${task.id}`}
                  className="flex cursor-pointer items-center"
                >
                  <Pencil className="mr-2 h-4 w-4" />
                  Edit task
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => handleStatusChange('Completed')}
                className="flex cursor-pointer items-center"
              >
                <CheckCircle className="mr-2 h-4 w-4" />
                Mark completed
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => handleStatusChange('In Progress')}
                className="flex cursor-pointer items-center"
              >
                <RotateCw className="mr-2 h-4 w-4" />
                Mark in-progress
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => handleStatusChange('Pending')}
                className="flex cursor-pointer items-center"
              >
                <Clock className="mr-2 h-4 w-4" />
                Mark pending
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={handleDelete}
                className="flex cursor-pointer items-center text-destructive focus:text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete task
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </CardHeader>

        <CardContent className="p-4 pt-0">
          <p className="line-clamp-3 text-sm text-muted-foreground">
            {task.description || 'No description provided'}
          </p>
        </CardContent>

        <CardFooter className="flex flex-col items-start space-y-3 border-t p-4">
          <div className="flex flex-wrap gap-1.5">
            <StatusBadge status={localStatus} size="sm" />
            <StatusUpdater
              taskId={task.id}
              currentStatus={localStatus}
              onUpdated={handleStatusChange}
            />
            <PriorityBadge priority={task.priority} size="sm" />
            {taskTags.length > 0 ? (
              taskTags.map((tag) => (
                <TagBadge key={tag.id} tag={tag} size="sm" />
              ))
            ) : (
              <span className="text-xs text-muted-foreground italic ml-1">
                No tags
              </span>
            )}
          </div>

          <div className="flex w-full flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground">
            <div>
              Created{' '}
              {task.createdAt
                ? formatDistanceToNow(new Date(task.createdAt), {
                    addSuffix: true,
                  })
                : 'Unknown'}
            </div>

            {task.dueDate && (
              <div
                className={`font-medium ${
                  new Date(task.dueDate) < new Date()
                    ? 'text-destructive'
                    : 'text-muted-foreground'
                }`}
              >
                Due{' '}
                {formatDistanceToNow(new Date(task.dueDate), {
                  addSuffix: true,
                })}
              </div>
            )}
          </div>
        </CardFooter>
      </Card>
    </motion.div>
  );
}
