'use client';

import { Task } from '@/types';
import { TaskCard } from '@/components/task-card';

interface TaskGridProps {
  tasks: Task[];
  onDelete(taskId: string): void;
  onStatusChange(taskId: string, status: Task['status']): void;
}

export function TaskGrid({ tasks, onDelete, onStatusChange }: TaskGridProps) {
  if (!Array.isArray(tasks) || tasks.length === 0) {
    return (
      <div className="flex min-h-[200px] flex-col items-center justify-center rounded-lg border border-dashed p-8 text-center">
        <h3 className="text-lg font-semibold">No tasks found</h3>
        <p className="mt-1 text-sm text-muted-foreground">
          Create a new task or adjust your filters to see results.
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3">
      {tasks.map((task, index) => (
        <TaskCard
          key={task.id}
          task={task}
          index={index}
          onDelete={() => onDelete(task.id)}
          onStatusChange={(status) => onStatusChange(task.id, status)}
        />
      ))}
    </div>
  );
}
