export const signInDefaultValues = {
  email: '',
  password: '',
};
export const signUpDefaultValues = {
  email: '',
  password: '',
  name: '',
  confirmPassword: '',
};
export const mapTask = (taskFromApi: any) => ({
  id: taskFromApi.id,
  title: taskFromApi.title,
  description: taskFromApi.description,
  priority: taskFromApi.priority,
  status: taskFromApi.status,
  dueDate: taskFromApi.due_date,
  createdAt: taskFromApi.created_at,
  updatedAt: taskFromApi.updated_at,
  userId: taskFromApi.user_id,
  tags: taskFromApi.tags,
});

export const mapComment = (commentFromApi: any) => ({
  id: commentFromApi.id,
  taskId: commentFromApi.task_id,
  userId: commentFromApi.user_id,
  content: commentFromApi.content,
  createdAt: new Date(commentFromApi.created_at),
  updatedAt: new Date(commentFromApi.updated_at),
});
