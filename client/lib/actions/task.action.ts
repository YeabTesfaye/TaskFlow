import { tasks } from '@/api';

export async function getTasks() {
  try {
    const response = await tasks.getAll();
    return response.tasks || [];
  } catch (error) {
    console.error('Error fetching tasks:', error);
    return [];
  }
}

export async function getTask(id: string) {
  try {
    return await tasks.getById(id);
  } catch (error) {
    console.error('Error fetching task:', error);
    return null;
  }
}