import React from 'react';
import styles from './TaskForm.module.css';

interface TaskFormProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (e: React.FormEvent) => void;
}

const TaskForm: React.FC<TaskFormProps> = ({ value, onChange, onSubmit }) => {
  return (
    <form onSubmit={onSubmit} className={styles.form}> 
      <input
        type="text"
        value={value}
        onChange={onChange}
        placeholder="Добавить задачу..."
        className={styles.input}
        required
      />
      <button type="submit" className={styles.button}>Добавить</button>
    </form>
  );
};

export default TaskForm;
