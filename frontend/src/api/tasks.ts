import { Task } from '../types/task';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || '/api';

export const getTasks = async (): Promise<Task[]> => {
  const response = await fetch(`${API_BASE_URL}/tasks`);
  if (!response.ok) {
    throw new Error('Failed to fetch tasks');
  }
  return response.json();
};

export const addTask = async (title: string): Promise<Task> => {
  const formData = new FormData();
  formData.append('title', title);
  
  const response = await fetch(`${API_BASE_URL}/add`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    throw new Error('Failed to add task');
  }
  return response.json();
};

export const completeTask = async (id: number): Promise<void> => {
  const formData = new FormData();
  formData.append('id', id.toString());
  
  const response = await fetch(`${API_BASE_URL}/done`, {
    method: 'POST',
    body: formData,
  });
  
  if (!response.ok) {
    throw new Error('Failed to complete task');
  }
};

export const deleteTask = async (id: number): Promise<void> => {
  const formData = new FormData();
  formData.append('id', id.toString());
  
  const response = await fetch(`${API_BASE_URL}/delete`, {
    method: 'POST',
    body: formData,
  });
  
  if (!response.ok) {
    throw new Error('Failed to delete task');
  }
};
