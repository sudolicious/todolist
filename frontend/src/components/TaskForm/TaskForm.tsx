import { Task } from '../../types/task';
import styles from './TaskForm.module.css';

interface TaskFormProps {
  value: string;
  onChange: (value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
}

const TaskForm = ({ value, onChange, onSubmit }: TaskFormProps) => {
  return (
    <form onSubmit={onSubmit} className={styles.form}>
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Enter a new task..."
        className={styles.input}
        required
      />
      <button type="submit" className={styles.button}>
        Add Task
      </button>
    </form>
  );
};

export default TaskForm;
