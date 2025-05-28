'use client';

import { Task } from '@/types/task';
import { TaskCard } from '@/components/task-card';

interface TaskGridProps {
  tasks: Task[];
}

export function TaskGrid({ tasks }: TaskGridProps) {
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
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {tasks.map((task, index) => (
        <TaskCard key={task.id} task={task} index={index} />
      ))}
    </div>
  );
}
