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

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewTaskTitle(e.target.value);
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>To Do List: Список задач</h1>
      <TaskForm 
        value={newTaskTitle}
        onChange={handleInputChange}
        onSubmit={handleSubmit}
      />
      {tasks && tasks.length > 0 ? (
        <TaskList 
          tasks={tasks} 
          onComplete={completeTask}
          onDelete={deleteTask}
        />
      ) : (
        <p>Нет задач. Добавьте первую!</p>
      )}
    </div>
  );
};

export default HomePage;
