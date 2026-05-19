<script setup lang="ts">
import {
  XMarkIcon,
  CheckCircleIcon,
  ClockIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'

defineEmits<{ close: [] }>()

const taskStore = useTaskStore()
</script>

<template>
  <Transition name="backdrop">
    <div
      class="fixed inset-0 bg-emerald-50/30 backdrop-blur-sm z-50 flex items-center justify-center p-4 dark:bg-slate-950/50"
      @click.self="$emit('close')"
    >

      <Transition name="modal">
        <div class="bg-emerald-50 border border-emerald-200 rounded-2xl w-full max-w-md shadow-2xl z-70 overflow-hidden dark:bg-slate-900 dark:border-slate-800 transition-colors">

          <div class="flex items-center justify-between px-5 py-4 border-b border-emerald-200 dark:border-slate-800 transition-colors">
            <h2 class="text-2sm font-semibold text-black tracking-wide dark:text-slate-100">Your statistics</h2>
            <button
              @click="$emit('close')"
              class="text-gray-900 hover:text-rose-600 hover:cursor-pointer transition-colors dark:text-slate-400 dark:hover:text-rose-400"
            >
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>

          <div class="px-5 py-5 space-y-5">

            <div class="flex items-center gap-5">
              <div class="relative w-16 h-16 flex-shrink-0">
                <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
                  <circle
                    cx="18" cy="18" r="16"
                    fill="none"
                    class="stroke-[oklch(92.5%_0.084_155.995)] dark:stroke-slate-700 transition-colors"
                    stroke-width="3"
                  />
                  <circle
                    cx="18" cy="18" r="16"
                    fill="none"
                    class="stroke-[oklch(69.6%_0.17_162.48)] dark:stroke-emerald-500 transition-colors"
                    stroke-width="3"
                    stroke-linecap="round"
                    :stroke-dasharray="`${taskStore.stats.rate} ${100 - taskStore.stats.rate}`"
                    stroke-dashoffset="0"
                  />
                </svg>
                <span class="absolute inset-0 flex items-center justify-center text-sm font-semibold text-black dark:text-slate-100">
                  {{ taskStore.stats.rate }}%
                </span>
              </div>

              <div>
                <p class="text-gray-900 font-semibold text-2sm dark:text-slate-100">
                  {{ taskStore.stats.done }}
                  <span class="text-gray-900 font-normal text-2sm dark:text-slate-400">/ {{ taskStore.stats.total }}</span>
                </p>
                <p class="text-gray-900 text-xs mt-1 dark:text-slate-400">tasks completed</p>
              </div>
            </div>

            <div class="border-t border-emerald-200 dark:border-slate-800 transition-colors"></div>

            <div class="grid grid-cols-3 gap-2">
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200 dark:bg-slate-800/50 dark:border-slate-700 transition-colors">
                <CheckCircleIcon class="w-4 h-4 text-emerald-600 mx-auto mb-1.5 dark:text-emerald-500" />
                <p class="text-black font-semibold text-lg dark:text-slate-100">{{ taskStore.stats.done }}</p>
                <p class="text-gray-800 text-xs mt-1 dark:text-slate-400">Done</p>
              </div>
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200 dark:bg-slate-800/50 dark:border-slate-700 transition-colors">
                <ClockIcon class="w-4 h-4 text-yellow-600 mx-auto mb-1.5 dark:text-yellow-500" />
                <p class="text-black font-semibold text-lg dark:text-slate-100">{{ taskStore.stats.active }}</p>
                <p class="text-gray-800 text-xs mt-1 dark:text-slate-400">Active</p>
              </div>
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200 dark:bg-slate-800/50 dark:border-slate-700 transition-colors">
                <ExclamationTriangleIcon class="w-4 h-4 text-rose-600 mx-auto mb-1.5 dark:text-rose-500" />
                <p class="text-black font-semibold text-lg dark:text-slate-100">{{ taskStore.stats.overdue }}</p>
                <p class="text-gray-800 text-xs mt-1 dark:text-slate-400">Overdue</p>
              </div>
            </div>

            <div class="space-y-2.5">
              <p class="text-xs text-gray-900 font-medium tracking-widest dark:text-slate-400">By priority</p>

              <div
                v-for="item in [
                  { label: 'High', value: taskStore.stats.byPriority.high, color: 'bg-red-600 dark:bg-red-500' },
                  { label: 'Medium', value: taskStore.stats.byPriority.medium, color: 'bg-yellow-600 dark:bg-yellow-500' },
                  { label: 'Low', value: taskStore.stats.byPriority.low, color: 'bg-gray-600 dark:bg-gray-500' },
                ]"
                :key="item.label"
                class="flex items-center gap-3"
              >
                <span class="text-xs text-gray-900 w-12 flex-shrink-0 dark:text-slate-300">{{ item.label }}</span>
                <div class="flex-1 h-2 bg-green-200/90 rounded-full overflow-hidden dark:bg-slate-700">
                  <div
                    :class="['h-full rounded-full transition-colors', item.color]"
                    :style="{
                      width: taskStore.stats.total > 0
                        ? `${Math.round((item.value / taskStore.stats.total) * 100)}%`
                        : '0%'
                    }"
                  ></div>
                </div>
                <span class="text-xs text-gray-900 w-4 text-right flex-shrink-0 dark:text-slate-300">{{ item.value }}</span>
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
</style>
