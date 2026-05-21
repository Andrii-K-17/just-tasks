<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  XMarkIcon,
  SparklesIcon,
  TrashIcon,
  CalendarIcon,
  PlusIcon,
  TagIcon,
} from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'
import { useCategoryStore } from '@/stores/useCategoryStore'
import { generateTasks } from '@/api/ai'

interface EditableTask {
  text: string
  deadline: string
  priority: 'low' | 'medium' | 'high'
}

const PRIORITY_NUM: Record<string, 1 | 2 | 3> = { low: 1, medium: 2, high: 3 }

const priorityOptions = [
  { value: 'low' as const,    label: 'Low',  color: 'text-gray-700 dark:text-gray-300' },
  { value: 'medium' as const, label: 'Med',  color: 'text-yellow-700 dark:text-yellow-500' },
  { value: 'high' as const,   label: 'High', color: 'text-red-700 dark:text-red-400' },
]

const props = defineProps<{ prompt: string }>()
const emit = defineEmits<{ close: [], saved: [] }>()

const taskStore = useTaskStore()
const categoryStore = useCategoryStore()

const loading = ref(true)
const saving = ref(false)
const error = ref('')
const categoryName = ref('')
const tasks = ref<EditableTask[]>([])

/**
 * Fetches AI-generated category and subtasks from the backend on mount.
 */
onMounted(async () => {
  try {
    const result = await generateTasks(props.prompt)
    categoryName.value = result.category
    tasks.value = result.tasks.map(t => ({
      text: t.text,
      deadline: t.deadline,
      priority: t.priority,
    }))
  } catch (e: any) {
    error.value = e.message || 'Failed to generate tasks'
  } finally {
    loading.value = false
  }
})

/**
 * Removes a task from the preview list by index.
 */
function removeTask(index: number) {
  tasks.value.splice(index, 1)
}

/**
 * Appends a blank task entry to the preview list.
 */
function addTask() {
  tasks.value.push({ text: '', deadline: '', priority: 'medium' })
}

/**
 * Validates inputs, creates the category and all non-empty tasks, then emits saved.
 */
async function save() {
  const validTasks = tasks.value.filter(t => t.text.trim())

  if (!categoryName.value.trim()) {
    error.value = 'Category name is required'
    return
  }
  if (validTasks.length === 0) {
    error.value = 'Add at least one task'
    return
  }

  saving.value = true
  error.value = ''

  try {
    const category = await categoryStore.add(categoryName.value.trim())

    const tasksToSave = [...validTasks].reverse()

    for (const task of tasksToSave) {
      await taskStore.add({
        task_text: task.text.trim(),
        priority: PRIORITY_NUM[task.priority] ?? 2,
        deadline: task.deadline || null,
        category_id: category.id,
      })
    }

    emit('saved')
  } catch (e: any) {
    error.value = e.message || 'Failed to save'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <Transition name="backdrop">
    <div
      class="fixed inset-0 bg-emerald-50/30 backdrop-blur-sm z-50 flex items-center justify-center p-4 dark:bg-slate-950/50"
      @click.self="emit('close')"
    >
      <Transition name="modal">
        <div class="flex flex-col bg-emerald-50 border border-emerald-200 rounded-2xl w-full max-w-2xl max-h-[calc(100vh-2rem)] shadow-2xl dark:bg-slate-900 dark:border-slate-800">

          <div class="flex items-center justify-between px-5 py-4 border-b border-emerald-200 dark:border-slate-800 transition-colors">
            <div class="flex items-center gap-2">
              <SparklesIcon class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
              <h2 class="text-2sm font-semibold text-black tracking-wide dark:text-slate-100">Generated Tasks Review</h2>
            </div>
            <button
              @click="emit('close')"
              class="text-gray-900 hover:text-rose-600 hover:cursor-pointer transition-colors dark:text-slate-400 dark:hover:text-rose-400"
            >
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>

          <div v-if="loading" class="px-5 py-12 flex flex-col items-center gap-3">
            <SparklesIcon class="w-8 h-8 text-emerald-500 animate-pulse dark:text-emerald-400" />
            <p class="text-sm text-gray-700 dark:text-slate-400">Generating your task plan...</p>
          </div>

          <div v-else-if="error && tasks.length === 0" class="px-5 py-10 flex flex-col items-center gap-3">
            <p class="text-sm text-rose-600 dark:text-rose-400">{{ error }}</p>
            <button
              @click="emit('close')"
              class="px-4 py-2 text-sm rounded-xl border border-emerald-200 text-emerald-700 hover:bg-emerald-100 transition-colors hover:cursor-pointer dark:border-slate-700 dark:text-emerald-400 dark:hover:bg-slate-800"
            >
              Close
            </button>
          </div>

          <div v-else class="flex-1 flex flex-col min-h-0 px-5 py-1.5 gap-2">

            <div class="flex items-center gap-2 flex-shrink-0">
              <TagIcon class="w-4 h-4 text-emerald-600 dark:text-emerald-400 flex-shrink-0" />
              <input
                v-model="categoryName"
                type="text"
                placeholder="Category name..."
                class="flex-1 bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-1.5 text-sm text-slate-900 placeholder-gray-500 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:placeholder-slate-400 dark:focus:border-emerald-500"
              />
            </div>

            <div class="border-t border-emerald-200 dark:border-slate-800 transition-colors flex-shrink-0"></div>

            <ul class="flex-1 overflow-y-auto space-y-2.5 pr-1 no-scrollbar">
              <li
                v-for="(task, index) in tasks"
                :key="index"
                class="flex flex-col gap-2 bg-white border border-emerald-200 rounded-xl p-3 dark:bg-slate-800/50 dark:border-slate-700 transition-colors"
              >
                <input
                  v-model="task.text"
                  type="text"
                  placeholder="Task description..."
                  class="w-full bg-transparent text-sm text-slate-900 placeholder-gray-500 focus:outline-none dark:text-slate-100 dark:placeholder-slate-400"
                />
                <div class="flex items-center gap-2 flex-wrap">

                  <div class="flex rounded-xl overflow-hidden shrink-0">
                    <button
                      v-for="opt in priorityOptions"
                      :key="opt.value"
                      type="button"
                      @click="task.priority = opt.value"
                      :class="[
                        'px-2.5 hover:cursor-pointer py-0.5 text-xs font-medium transition-colors',
                        opt.color,
                        task.priority === opt.value
                          ? 'bg-emerald-300/80 rounded-xl border border-emerald-200 dark:bg-emerald-900/60 dark:border-emerald-700'
                          : 'rounded-xl border border-transparent hover:border-emerald-200 dark:hover:border-slate-700'
                      ]"
                    >
                      {{ opt.label }}
                    </button>
                  </div>

                  <div class="relative">
                    <CalendarIcon class="w-3.5 h-3.5 absolute right-3 top-1/2 -translate-y-1/2 text-emerald-600 dark:text-emerald-500 pointer-events-none" />
                    <input
                      v-model="task.deadline"
                      type="date"
                      class="hover:cursor-pointer min-w-0 max-w-[130px] bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-0.5 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:focus:border-emerald-500"
                    />
                  </div>

                  <button
                    type="button"
                    @click="removeTask(index)"
                    class="ml-auto text-gray-400 hover:text-rose-500 hover:cursor-pointer transition-colors dark:text-slate-500 dark:hover:text-rose-400"
                    aria-label="Remove task"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </li>
            </ul>

            <div class="flex flex-col gap-3 pt-2 bg-emerald-50 dark:bg-slate-900 flex-shrink-0">
              <button
                type="button"
                @click="addTask"
                class="w-full py-1.5 border border-dashed border-emerald-300 rounded-xl text-xs text-emerald-600 hover:bg-emerald-100 transition-colors flex items-center justify-center gap-1 hover:cursor-pointer dark:border-emerald-700 dark:text-emerald-400 dark:hover:bg-slate-800/50"
              >
                <PlusIcon class="w-3.5 h-3.5" /> Add task
              </button>

              <p v-if="error" class="text-rose-600 text-xs dark:text-rose-400">{{ error }}</p>

              <div class="flex gap-2">
                <button
                  type="button"
                  @click="emit('close')"
                  class="flex-1 py-2 text-sm rounded-xl border border-emerald-200 text-gray-700 hover:bg-emerald-100 transition-colors hover:cursor-pointer dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  @click="save"
                  :disabled="saving"
                  class="flex-1 py-2 text-sm font-medium rounded-xl bg-emerald-600 text-white hover:bg-emerald-700 transition-colors hover:cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed dark:bg-emerald-600 dark:hover:bg-emerald-500"
                >
                  {{ saving ? 'Saving...' : 'Save all' }}
                </button>
              </div>
            </div>

          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
.backdrop-enter-active, .backdrop-leave-active {
  transition: opacity .2s ease;
}
.backdrop-enter-from, .backdrop-leave-to {
  opacity: 0;
}

.modal-enter-active {
  transition: all .2s ease;
}
.modal-leave-active {
  transition: all .15s ease;
}
.modal-enter-from {
  opacity: 0;
  transform: scale(.96) translateY(8px);
}
.modal-leave-to {
  opacity: 0;
  transform: scale(.96);
}

input[type="date"]::-webkit-calendar-picker-indicator {
  opacity: 0;
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
