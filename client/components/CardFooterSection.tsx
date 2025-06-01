import { CardFooter } from '@/components/ui/card';
import { StatusBadge } from '@/components/ui/status-badge';
import { PriorityBadge } from '@/components/ui/priority-badge';
import { TagBadge } from '@/components/ui/tag-badge';
import { StatusUpdater } from './status-updater';
import { formatDistanceToNow } from 'date-fns';
import { Task, Tag } from '@/types';

interface Props {
  task: Task;
  status: Task['status'];
  tags: Tag[];
  onStatusChange(status: Task['status']): void;
}

export function CardFooterSection({
  task,
  status,
  tags,
  onStatusChange,
}: Props) {
  return (
    <CardFooter className="flex flex-col items-start space-y-3 border-t p-4">
      <div className="flex flex-wrap gap-1.5">
        <StatusBadge status={status} size="sm" />
        <StatusUpdater
          taskId={task.id}
          currentStatus={status}
          onUpdated={onStatusChange}
        />
        <PriorityBadge priority={task.priority} size="sm" />
        {tags.length > 0 ? (
          tags.map((tag) => <TagBadge key={tag.id} tag={tag} size="sm" />)
        ) : (
          <span className="text-xs italic text-muted-foreground ml-1">
            No tags
          </span>
        )}
      </div>
      <div className="flex w-full flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground">
        <div>
          Created{' '}
          {task.createdAt
            ? formatDistanceToNow(new Date(task.createdAt), { addSuffix: true })
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
  );
}
