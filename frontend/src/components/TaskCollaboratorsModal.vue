<script setup lang="ts">
import { ref, computed } from 'vue'
import { XMarkIcon, UserIcon, UserPlusIcon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'
import { useAuthStore } from '@/stores/useAuthStore'
import type { Task } from '@/types'

const props = defineProps<{ task: Task }>()
const emit = defineEmits<{ (e: 'close'): void }>()

const taskStore = useTaskStore()
const authStore = useAuthStore()

const isOwner = computed(() => authStore.user?.username === props.task.owner_name)
const showCollabMenu = ref(true)
const newCollabUsername = ref('')

async function handleAddCollaborator() {
  const username = newCollabUsername.value.trim()
  if (username) {
    try {
      await taskStore.addTaskCollaborator(props.task.id, username)
      newCollabUsername.value = ''
    } catch (e: any) {
      alert(e.message || 'Failed to add collaborator')
    }
  }
}

async function handleRemoveCollaborator(collabId: number) {
  if (confirm('Remove this collaborator?')) {
    await taskStore.removeTaskCollaborator(props.task.id, collabId)
  }
}
</script>

<template>
  <Transition name="backdrop">
    <div
      v-if="showCollabMenu && isOwner"
      class="fixed inset-0 bg-emerald-50/30 backdrop-blur-sm z-50 flex items-center justify-center p-4 dark:bg-slate-950/50"
      @click.self="emit('close')"
    >
      <Transition name="modal">
        <div class="bg-white border border-emerald-200 rounded-2xl w-full max-w-sm shadow-2xl z-70 overflow-hidden dark:bg-slate-900 dark:border-slate-800 transition-colors">
          <div class="flex items-center justify-between px-5 py-4 border-b border-emerald-200 dark:border-slate-800">
            <h2 class="text-2sm font-semibold text-black tracking-wide dark:text-slate-100">Share this task</h2>
            <button
              @click="emit('close')"
              class="text-gray-900 cursor-pointer hover:text-rose-600 transition-colors dark:text-slate-400 dark:hover:text-rose-400"
            >
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>

          <div class="px-5 py-5 space-y-4">
            <div class="flex items-center gap-2">
              <input
                v-model="newCollabUsername"
                @keydown.enter="handleAddCollaborator"
                type="text"
                placeholder="Username..."
                class="flex-1 px-3 py-1.5 text-sm text-gray-900 border border-gray-200 rounded-lg focus:outline-none focus:border-emerald-500 dark:bg-slate-900 dark:border-slate-600 dark:text-gray-100"
              />
              <button
                @click="handleAddCollaborator"
                class="p-2 cursor-pointer transition-all bg-emerald-100 text-emerald-700 rounded-lg hover:bg-emerald-200 dark:bg-emerald-600 dark:hover:bg-emerald-500 dark:text-emerald-50"
              >
                <UserPlusIcon class="w-4.5 h-4.5" />
              </button>
            </div>

            <div v-if="props.task.collaborators && props.task.collaborators.length > 0" class="space-y-2">
              <h4 class="text-[10px] uppercase text-gray-500 mb-1 font-semibold">Shared with</h4>
              <ul class="space-y-1 max-h-40 overflow-y-auto no-scrollbar">
                <li
                  v-for="collab in props.task.collaborators"
                  :key="collab.id"
                  class="flex items-center justify-between text-sm p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700/50"
                >
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-700 dark:bg-slate-800 dark:text-emerald-400 text-xs font-semibold">
                        <UserIcon class="w-4.5 h-4.5" />
                    </div>
                    <span class="truncate max-w-[130px] text-slate-900 dark:text-slate-300">{{ collab.username }}</span>
                  </div>
                  <button
                    @click="handleRemoveCollaborator(collab.id)"
                    class="text-gray-600 dark:text-gray-400 cursor-pointer hover:text-rose-500"
                  >
                    <XMarkIcon class="w-4.5 h-4.5" />
                  </button>
                </li>
              </ul>
            </div>
            <div v-else class="text-xs text-gray-500 dark:text-slate-400">No collaborators yet</div>
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
</style>
