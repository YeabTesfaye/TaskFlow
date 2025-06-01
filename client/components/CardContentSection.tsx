import { CardContent } from '@/components/ui/card';
import { Task } from '@/types';

export function CardContentSection({ task }: { task: Task }) {
  return (
    <CardContent className="p-4 pt-0">
      <p className="line-clamp-3 text-sm text-muted-foreground">
        {task.description || 'No description provided'}
      </p>
    </CardContent>
  );
}
