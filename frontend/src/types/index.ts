export interface TaskCollaborator {
  id: number
  username: string
}

export interface Task {
  id: number
  user_id: number
  task_text: string
  priority: 1 | 2 | 3
  deadline: string | null
  is_completed: boolean
  category_id: number | null
  created_at: string
  owner_name?: string
  collaborators?: TaskCollaborator[]
}

export interface User {
  id: number
  username: string
}

export interface Category {
  id: number
  name: string
}
