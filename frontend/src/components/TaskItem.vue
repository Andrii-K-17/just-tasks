<script setup lang="ts">
import { ref, computed } from 'vue'
import { TrashIcon, CheckIcon, TagIcon, CalendarIcon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'
import type { Task } from '@/types'
import { useTextareaAutosize, onClickOutside } from '@vueuse/core'
import { useCategoryStore } from '@/stores/useCategoryStore'

const props = defineProps<{ task: Task }>()
const taskStore = useTaskStore()
const categoryStore = useCategoryStore()

const isEditing = ref(false)
const editValue = ref('')

const isEditingDeadline = ref(false)
const editDeadlineValue = ref('')

const editCategoryValue = ref<number | null>(null)
const isEditingCategory = ref(false)

const deadlineCard = ref(null)
const categoryCard = ref(null)

const textareaRef = ref<HTMLTextAreaElement | null>(null)
useTextareaAutosize({ element: textareaRef })

const PRIORITY: Record<number, string> = {
  1: 'bg-gray-500',
  2: 'bg-yellow-500',
  3: 'bg-red-500',
}

const today = new Date().toISOString().slice(0, 10)
const isOverdue = computed(
  () => props.task.deadline && !props.task.is_completed && props.task.deadline < today
)

/**
 * Initializes the text editing mode with the current task text.
 */
function startEdit() {
  editValue.value = props.task.task_text
  isEditing.value = true
}

/**
 * Persists the updated task text via the store if the content has changed.
 */
async function saveEdit() {
  const trimmed = editValue.value.trim()
  if (trimmed && trimmed !== props.task.task_text) {
    await taskStore.editText(props.task.id, trimmed)
  }
  isEditing.value = false
}

/**
 * Initializes the deadline editing mode with the current deadline value.
 */
function startEditDeadline() {
  editDeadlineValue.value = props.task.deadline?? ''
  isEditingDeadline.value = true
}

/**
 * Persists the updated task deadline via the store if the date has changed.
 */
async function saveEditDeadline() {
  const value = editDeadlineValue.value
  if (value !== props.task.deadline) {
    await taskStore.editDeadline(props.task.id, value)
  }
  isEditingDeadline.value = false
}

/**
 * Initializes the category editing mode with the current category ID.
 */
function startEditCategory() {
  editCategoryValue.value = props.task.category_id ?? null
  isEditingCategory.value = true
}

/**
 * Persists the updated task category via the store if it has changed.
 */
async function saveEditCategory() {
  const value = editCategoryValue.value
  if (value !== props.task.category_id) {
    await taskStore.editCategory(props.task.id, value)
  }
  isEditingCategory.value = false
}

/**
 * Handles keyboard shortcuts to save changes on Enter or cancel editing on Escape.
 */
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    if (isEditing.value) saveEdit()
    if (isEditingDeadline.value) saveEditDeadline()
  }
  if (e.key === 'Escape') {
    isEditing.value = false
    isEditingDeadline.value = false
  }
}

onClickOutside(textareaRef, () => {
  saveEdit()
})

onClickOutside(deadlineCard, () => {
  saveEditDeadline()
})

onClickOutside(categoryCard, () => {
  saveEditCategory()
})
</script>

<template>
  <li
    class="flex items-center gap-3 bg-white border border-emerald-200 rounded-2xl px-4 py-3 group hover:border-emerald-500 transition-colors dark:bg-slate-900 dark:border-slate-800 dark:hover:border-emerald-400"
    @click="startEdit"
  >
    <button
      @click.stop="taskStore.toggle(task.id)"
      :class="[
        'w-5 h-5 hover:cursor-pointer rounded-full border-1 flex items-center justify-center flex-shrink-0 transition-all',
        task.is_completed
          ? 'bg-emerald-600 border-emerald-600 dark:bg-emerald-500 dark:border-emerald-500'
          : 'border-emerald-900 hover:border-emerald-500 dark:border-slate-500 dark:hover:border-emerald-400'
      ]"
      :aria-label="task.is_completed ? 'Mark incomplete' : 'Mark complete'"
    >
      <CheckIcon v-if="task.is_completed" class="w-3 h-3 text-white dark:text-slate-900" />
    </button>

    <span
      :class="['w-[0.55%] h-5 rounded-full flex-shrink-0', PRIORITY[task.priority]]"
    ></span>

    <div class="flex-1 min-w-0">
      <textarea
        v-if="isEditing"
        v-model="editValue"
        ref="textareaRef"
        @blur="saveEdit"
        @keydown="onKeydown"
        class="w-full bg-slate-50 rounded-lg px-2 py-0.5 text-sm text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-300 resize-none overflow-hidden dark:bg-slate-800 dark:text-slate-100 dark:focus:ring-emerald-500"
      ></textarea>
      <span
        v-else
        @click="startEdit"
        :class="[
          'text-sm cursor-text block truncate select-none',
          task.is_completed ? 'line-through text-gray-700 dark:text-slate-500' : 'text-gray-900 dark:text-slate-100'
        ]"
        :title="task.task_text"
      >
        {{ task.task_text }}
      </span>

      <div class="flex items-center gap-1">
        <div
          v-if="isEditingDeadline"
          ref="deadlineCard"
          class="relative"
        >
          <CalendarIcon class="w-3.5 h-3.5 absolute right-3 top-1/2 -translate-y-1/2 text-emerald-600 dark:text-emerald-500 pointer-events-none" />
          <input
            v-model="editDeadlineValue"
            type="date"
            @blur="saveEditDeadline"
            @keydown="onKeydown"
            class="hover:cursor-pointer flex-1 min-w-0 max-w-[140px] bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:focus:border-emerald-500"
          >
        </div>
        <div
          v-else-if="task.deadline && !isEditingDeadline"
          class = "relative"
        >
          <CalendarIcon class="w-3.5 h-3.5 absolute top-1/2 -translate-y-1/2 text-emerald-600 dark:text-emerald-500 pointer-events-none" />
          <span
            @click="startEditDeadline"
            :class="['pl-5 text-xs block mt-0.5 mr-2', isOverdue ? 'text-rose-600 dark:text-rose-400' : 'text-gray-700 dark:text-slate-400']"
          >
            {{ task.deadline }}
          </span>
        </div>
        <button
          v-else-if="!task.deadline && !isEditingDeadline && isEditing"
          @click="startEditDeadline"
          class="flex-1 hover:bg-emerald-50 dark:hover:bg-slate-700 hover:cursor-pointer min-w-0 bg-emerald-50/20 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-600 dark:text-slate-100 dark:focus:border-emerald-500"
          aria-label="Set deadline"
        >
          Set deadline
        </button>
        
        <div
          v-if="isEditingCategory"
          ref="categoryCard"
          class="relative flex-1 min-w-[120px]"
        >
          <TagIcon class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-emerald-600 dark:text-emerald-500 pointer-events-none" />
          <select
            v-model="editCategoryValue"
            @blur="saveEditCategory"
            class="hover:cursor-pointer appearance-none bg-emerald-50/30 border border-emerald-200 rounded-xl pl-8 pr-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:focus:border-emerald-500"
          >
            <option :value="null">No category</option>
            <option
              v-for="category in categoryStore.categories"
              :key="category.id"
              :value="category.id"
            >
              {{ category.name }}
            </option>
          </select>
        </div>
        <div v-else-if="task.category_id && categoryStore.getById(task.category_id)" class="flex items-center gap-1">
          <TagIcon class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-500 pointer-events-none" />  
          <span
            @click="startEditCategory"
            class="text-xs text-gray-700 dark:text-slate-400"
          >
            {{ categoryStore.getById(task.category_id)?.name }}
          </span>
        </div>
        <button
          v-else-if="(!task.category_id || !categoryStore.getById(task.category_id)) && isEditing && categoryStore.hasCategories"
          @click="startEditCategory"
          class="flex-1 hover:cursor-pointer hover:bg-emerald-50 dark:hover:bg-slate-700 min-w-0 bg-emerald-50/20 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-600 dark:text-slate-100 dark:focus:border-emerald-500"
          aria-label="Set category"
        >
          Set category
        </button>
      </div>

    </div>

    <button
      @click.stop="taskStore.remove(task.id)"
      class="opacity-0 hover:cursor-pointer group-hover:opacity-100 text-gray-900 hover:text-rose-600 transition-all flex-shrink-0 dark:text-slate-400 dark:hover:text-rose-400"
      aria-label="Delete task"
    >
      <TrashIcon class="w-5 h-5" />
    </button>
  </li>
</template>

<style scoped>
input[type="date"]::-webkit-calendar-picker-indicator {
  opacity: 0;
}
</style>
