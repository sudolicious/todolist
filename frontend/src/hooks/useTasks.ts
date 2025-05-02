import { useState, useEffect, useCallback } from 'react';
import { Task } from '../types/task';
import { getTasks, addTask, completeTask, deleteTask } from '../api/tasks';

export const useTasks = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchTasks = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getTasks();
      setTasks(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const handleAddTask = async (title: string) => {
    try {
      const newTask = await addTask(title);
      setTasks(prev => [...prev, newTask]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add task');
    }
  };

  const handleCompleteTask = async (id: number) => {
    try {
      await completeTask(id);
      setTasks(prev => prev.map(task => 
        task.id === id ? { ...task, done: true } : task
      ));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to complete task');
    }
  };

  const handleDeleteTask = async (id: number) => {
    try {
      await deleteTask(id);
      setTasks(prev => prev.filter(task => task.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete task');
    }
  };

  return {
    tasks,
    loading,
    error,
    addTask: handleAddTask,
    completeTask: handleCompleteTask,
    deleteTask: handleDeleteTask,
    refreshTasks: fetchTasks,
  };
};
