import { Task } from '../../types/task';
import TaskItem from '../TaskItem/TaskItem';
import styles from './TaskList.module.css';

interface TaskListProps {
  tasks: Task[];
  onComplete: (id: number) => void;
  onDelete: (id: number) => void;
}

const TaskList = ({ tasks, onComplete, onDelete }: TaskListProps) => {
  return (
    <ul className={styles.taskList}>
      {tasks.map(task => (
        <TaskItem
          key={task.id}
          task={task}
          onComplete={onComplete}
          onDelete={onDelete}
        />
      ))}
    </ul>
  );
};

export default TaskList;
