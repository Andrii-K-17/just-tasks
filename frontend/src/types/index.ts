export interface Task {
  id: number
  user_id: number
  task_text: string
  priority: 1 | 2 | 3
  deadline: string | null
  is_completed: boolean
  created_at: string
}

export interface User {
  id: number
  username: string
}