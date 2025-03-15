export type WeAll = 'Dasha' | 'Kirill' | 'IPA' | 'Shiki';

export interface TaskCreate {
  name: string;
  doer: WeAll;
  description?: string;
  do_before: string; // ISO 8601 format
  repeatable: boolean;
  created_by: WeAll;
  range?: number | null;
}

export interface Task extends TaskCreate {
  id: string;
  done: boolean;
}

export interface TodosResponse {
  tasks: Task[];
}

export interface ErrorResponse {
  error: string;
}