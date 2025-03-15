import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useParams } from 'react-router-dom';
import type { Task, TodosResponse, ErrorResponse, WeAll } from './types';

export default function App() {
  return (
    <Router>
      <nav style={{ padding: '1rem', backgroundColor: '#f0f0f0' }}>
        <Link to="/" style={{ marginRight: '1rem' }}>All Tasks</Link>
        <Link to="/create">Create Task</Link>
      </nav>

      <Routes>
        <Route path="/" element={<TodoList />} />
        <Route path="/create" element={<CreateTask />} />
        <Route path="/todo/:todoId" element={<TaskDetails />} />
      </Routes>
    </Router>
  );
}

function TodoList() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch('/api/v1/todo')
      .then(async (response) => {
        if (!response.ok) {
          const error: ErrorResponse = await response.json();
          throw new Error(error.error);
        }
        const data: TodosResponse = await response.json();
        setTasks(data.tasks);
      })
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div style={{ padding: '2rem' }}>
      <h2>All Tasks</h2>
      <div style={{ display: 'grid', gap: '1rem' }}>
        {tasks.map(task => (
          <div key={task.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '4px' }}>
            <h3 style={{ marginTop: 0 }}>
              <Link to={`/todo/${task.id}`} style={{ textDecoration: 'none', color: '#333' }}>
                {task.name}
              </Link>
            </h3>
            <p>Assignee: {task.doer}</p>
            <p>Due: {new Date(task.do_before).toLocaleString()}</p>
            <p>Status: {task.done ? '✅ Completed' : '🟡 Pending'}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

function CreateTask() {
  const [formData, setFormData] = useState<Omit<Task, 'id' | 'done'>>({
    name: '',
    doer: 'Dasha',
    description: '',
    do_before: '',
    repeatable: false,
    created_by: 'Dasha',
    range: null,
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const weAllOptions: WeAll[] = ['Dasha', 'Kirill', 'IPA', 'Shiki'];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Format the date before sending
    const formattedData = {
      ...formData,
      do_before: new Date(formData.do_before).toISOString(),
    };

    try {
      const response = await fetch('/api/v1/todo', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formattedData),
      });

      if (!response.ok) {
        const error: ErrorResponse = await response.json();
        throw new Error(error.error);
      }

      const data = await response.json();
      setMessage(`Task created successfully! ID: ${data.task_id}`);
      setError('');
      setFormData({
        name: '',
        doer: 'Dasha',
        description: '',
        do_before: '',
        repeatable: false,
        created_by: 'Dasha',
        range: null,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create task');
      setMessage('');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  return (
    <div style={{ padding: '2rem', maxWidth: '600px', margin: '0 auto' }}>
      <h2>Create New Task</h2>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: '1rem' }}>
        <div>
          <label>Task Name *</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            style={{ width: '100%' }}
          />
        </div>

        <div>
          <label>Assignee *</label>
          <select
            name="doer"
            value={formData.doer}
            onChange={handleChange}
            style={{ width: '100%' }}
          >
            {weAllOptions.map(option => (
              <option key={option} value={option}>{option}</option>
            ))}
          </select>
        </div>

        <div>
          <label>Description</label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            style={{ width: '100%', minHeight: '100px' }}
          />
        </div>

        <div>
          <label>Due Date *</label>
          <input
            type="datetime-local"
            name="do_before"
            value={formData.do_before}
            onChange={handleChange}
            required
            style={{ width: '100%' }}
          />
        </div>

        <div>
          <label>
            <input
              type="checkbox"
              name="repeatable"
              checked={formData.repeatable}
              onChange={handleChange}
            />
            Repeatable
          </label>
        </div>

        <div>
          <label>Created By *</label>
          <select
            name="created_by"
            value={formData.created_by}
            onChange={handleChange}
            style={{ width: '100%' }}
          >
            {weAllOptions.map(option => (
              <option key={option} value={option}>{option}</option>
            ))}
          </select>
        </div>

        {formData.repeatable && (
          <div>
            <label>Repeat Interval (days)</label>
            <input
              type="number"
              name="range"
              value={formData.range ?? ''}
              onChange={handleChange}
              min="1"
              style={{ width: '100%' }}
            />
          </div>
        )}

        {message && <div style={{ color: 'green' }}>{message}</div>}
        {error && <div style={{ color: 'red' }}>{error}</div>}

        <button
          type="submit"
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
          }}
        >
          Create Task
        </button>
      </form>
    </div>
  );
}

function TaskDetails() {
  const { todoId } = useParams<{ todoId: string }>();
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!todoId) return;

    fetch(`/api/v1/todo/${todoId}`)
      .then(async (response) => {
        if (!response.ok) {
          if (response.status === 404) throw new Error('Task not found');
          const error: ErrorResponse = await response.json();
          throw new Error(error.error);
        }
        const data: Task = await response.json();
        setTask(data);
      })
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, [todoId]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!task) return <div>Task not found</div>;

  return (
    <div style={{ padding: '2rem', maxWidth: '600px', margin: '0 auto' }}>
      <h2>{task.name}</h2>
      <div style={{ marginTop: '1rem', lineHeight: '1.6' }}>
        <p><strong>Assignee:</strong> {task.doer}</p>
        <p><strong>Description:</strong> {task.description || 'None'}</p>
        <p><strong>Due Date:</strong> {new Date(task.do_before).toLocaleString()}</p>
        <p><strong>Repeatable:</strong> {task.repeatable ? 'Yes' : 'No'}</p>
        <p><strong>Created By:</strong> {task.created_by}</p>
        <p><strong>Status:</strong> {task.done ? 'Completed' : 'Pending'}</p>
        {task.range && <p><strong>Repeat Interval:</strong> Every {task.range} days</p>}
      </div>
    </div>
  );
}