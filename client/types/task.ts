export type Priority = 'low' | 'medium' | 'high' | 'urgent';

export type TaskStatus = 'pending' | 'in-progress' | 'completed';

export interface Task {
  id: string;
  title: string;
  description: string;
  priority: Priority;
  status: TaskStatus;
  dueDate: Date | null;
  createdAt: Date;
  tags: string[];
}

export type TaskFormValues = Omit<Task, 'id' | 'createdAt'>;

export interface Tag {
  id: string;
  name: string;
  color: string;
}