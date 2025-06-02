import axios from 'axios';
import { Status } from './types';
import { mapComment, mapTask } from './lib/constant/inde';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to attach token
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Auth endpoints
export const auth = {
  login: async (email: string, password: string) => {
    const response = await api.post('/login', { email, password });
    return response.data;
  },
  signup: async (email: string, password: string, name: string) => {
    const response = await api.post('/signup', { email, password, name });
    return response.data;
  },
};

export const user = {
  getProfile: async () => {
    const response = await api.get('/users/me');
    return response.data;
  },
  updateProfile: async (data: any) => {
    const response = await api.patch('/users/me', data);
    return response.data;
  },
  changePassword: async (currentPassword: string, newPassword: string) => {
    const response = await api.post('/users/password', {
      current_password: currentPassword,
      new_password: newPassword,
    });
    return response.data;
  },
  deleteAccount: async () => {
    const response = await api.delete('/users/me');
    return response.data;
  },
  uploadProfilePicture: async (file: File) => {
    const formData = new FormData();
    formData.append('profile_picture', file, file.name);

    const response = await api.post('/users/profile-picture', formData, {
      headers: {
        'Content-Type': null,
      },
    });
    return response.data;
  },
  updatePreferences: async (preferences: {
    timezone?: string;
    emailNotifications: boolean;
    pushNotifications: boolean;
    dailyDigest: boolean;
  }) => {
    const payload = {
      timezone: preferences.timezone || 'UTC', 
      email_notifications: preferences.emailNotifications,
      push_notifications: preferences.pushNotifications,
      daily_digest: preferences.dailyDigest,
    };

    const response = await api.put('/users/preferences', payload);
    return response.data;
  },
};

export const tasks = {
  create: async (data: any) => {
    // Create a new object with only the required fields
    const taskData = {
      title: data.title,
      description: data.description,
      priority: data.priority,
      status: 'Pending', // Start with Pending only
      due_date: data.dueDate,
      tags: data.tags,
    };

    const response = await api.post('/tasks', taskData);
    return response.data;
  },
  getAll: async () => {
    const response = await api.get('/tasks');
    return {
      ...response.data,
      tasks: response.data.tasks.map(mapTask),
    };
  },
  getById: async (id: string) => {
    const response = await api.get(`/tasks/${id}`);
    return mapTask(response.data);
  },
  update: async (id: string, data: any) => {
    // Ensure status has correct format
    const formattedData = {
      ...data,
      status:
        data.status === 'In Progress'
          ? 'In Progress'
          : data.status.charAt(0).toUpperCase() + data.status.slice(1),
    };
    const response = await api.put(`/tasks/${id}`, formattedData);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/tasks/${id}`);
    return response.data;
  },
  updateStatus: async (id: string, data: { status: Status }) => {
    const response = await api.patch(`/tasks/${id}/status`, data);
    return response.data;
  },
  addCollaborator: async (taskId: string, collaboratorId: string) => {
    const response = await api.post(`/tasks/collaborators/add`, {
      task_id: taskId,
      collaborator_id: collaboratorId,
    });
    return response.data;
  },
  removeCollaborator: async (taskId: string, collaboratorId: string) => {
    const response = await api.post(`/tasks/collaborators/remove`, {
      task_id: taskId,
      collaborator_id: collaboratorId,
    });
    return response.data;
  },
};

export const tags = {
  getAll: async () => {
    const response = await api.get('/tags');
    return response.data;
  },
};

export const comments = {
  create: async (taskId: string, data: { content: string }) => {
    const response = await api.post(`/tasks/${taskId}/comments`, data);
    return mapComment(response.data);
  },
  getAll: async (taskId: string) => {
    const response = await api.get(`/tasks/${taskId}/comments`);
    const rawData = response.data;

    if (!rawData.comments || !Array.isArray(rawData.comments)) return [];

    return rawData.comments.map(mapComment);
  },
  update: async (
    taskId: string,
    commentId: string,
    data: { content: string },
  ) => {
    const response = await api.put(
      `/tasks/${taskId}/comments/${commentId}`,
      data,
    );
    return mapComment(response.data);
  },
  delete: async (taskId: string, commentId: string) => {
    const response = await api.delete(`/tasks/${taskId}/comments/${commentId}`);
    return response.data;
  },
};

export const statistics = {
  getTaskStatistics: async () => {
    const response = await api.get('/tasks/statistics');
    return response.data;
  },
  getTaskStatisticsByCategory: async () => {
    const response = await api.get('/tasks/statistics');
    return response.data;
  },
};
export default api;
