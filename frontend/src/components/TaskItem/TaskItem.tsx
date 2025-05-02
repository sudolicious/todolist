import { Task } from '../../types/task';
import styles from './TaskItem.module.css';

interface TaskItemProps {
  task: Task;
  onComplete: (id: number) => void;
  onDelete: (id: number) => void;
}

const TaskItem = ({ task, onComplete, onDelete }: TaskItemProps) => {
  return (
    <li className={styles.taskItem}>
      <div className={styles.taskContent}>
        <input
          type="checkbox"
          checked={task.done}
          onChange={() => onComplete(task.id)}
          className={styles.checkbox}
        />
        <span 
          className={`${styles.taskTitle} ${task.done ? styles.completed : ''}`}
        >
          {task.title}
        </span>
      </div>
      <button
        onClick={() => onDelete(task.id)}
        className={styles.deleteButton}
        aria-label="Delete task"
      >
        ×
      </button>
    </li>
  );
};

export default TaskItem;
