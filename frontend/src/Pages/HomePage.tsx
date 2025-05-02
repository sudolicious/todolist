import { useTasks } from '../hooks/useTasks';
import React, { useState } from 'react';
import TaskForm from '../components/TaskForm/TaskForm';
import TaskList from '../components/TaskList/TaskList';

const HomePage = () => {
  const { tasks, loading, error, addTask, completeTask, deleteTask } = useTasks();
  const [newTaskTitle, setNewTaskTitle] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newTaskTitle.trim()) {
      addTask(newTaskTitle.trim());
      setNewTaskTitle('');
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>To Do List</h1>
      <TaskForm 
        value={newTaskTitle}
        onChange={setNewTaskTitle}
        onSubmit={handleSubmit}
      />
      <TaskList 
        tasks={tasks}
        onComplete={completeTask}
        onDelete={deleteTask}
      />
    </div>
  );
};

export default HomePage;
