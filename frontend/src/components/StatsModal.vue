<script setup lang="ts">
import { XMarkIcon,
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
      class="fixed inset-0 bg-emerald-50/30 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      @click.self="$emit('close')"
    >

      <Transition name="modal">
        <div class="bg-emerald-50 border border-emerald-200 rounded-2xl w-full max-w-md shadow-2xl z-70 overflow-hidden">

          <div class="flex items-center justify-between px-5 py-4 border-b border-emerald-200">
            <h2 class="text-2sm font-semibold text-black tracking-wide">Your statistics</h2>
            <button
              @click="$emit('close')"
              class="text-gray-900 hover:text-rose-600 hover:cursor-pointer transition-colors"
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
                    stroke="oklch(92.5% 0.084 155.995)"
                    stroke-width="3"
                  />
                  <circle
                    cx="18" cy="18" r="16"
                    fill="none"
                    stroke="oklch(69.6% 0.17 162.48)"
                    stroke-width="3"
                    stroke-linecap="round"
                    :stroke-dasharray="`${taskStore.stats.rate} ${100 - taskStore.stats.rate}`"
                    stroke-dashoffset="0"
                  />
                </svg>
                <span class="absolute inset-0 flex items-center justify-center text-sm font-semibold text-black">
                  {{ taskStore.stats.rate }}%
                </span>
              </div>

              <div>
                <p class="text-gray-900 font-semibold text-2sm">
                  {{ taskStore.stats.done }}
                  <span class="text-gray-900 font-normal text-2sm">/ {{ taskStore.stats.total }}</span>
                </p>
                <p class="text-gray-900 text-xs mt-1">tasks completed</p>
              </div>
            </div>

            <div class="border-t border-emerald-200"></div>

            <div class="grid grid-cols-3 gap-2">
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200">
                <CheckCircleIcon class="w-4 h-4 text-emerald-600 mx-auto mb-1.5" />
                <p class="text-black font-semibold text-lg=">{{ taskStore.stats.done }}</p>
                <p class="text-gray-800 text-xs mt-1">Done</p>
              </div>
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200">
                <ClockIcon class="w-4 h-4 text-yellow-600 mx-auto mb-1.5" />
                <p class="text-black font-semibold text-lg">{{ taskStore.stats.active }}</p>
                <p class="text-gray-800 text-xs mt-1">Active</p>
              </div>
              <div class="bg-emerald-200/30 rounded-xl p-3 text-center border border-emerald-200">
                <ExclamationTriangleIcon class="w-4 h-4 text-rose-600 mx-auto mb-1.5" />
                <p class="text-black font-semibold text-lg">{{ taskStore.stats.overdue }}</p>
                <p class="text-gray-800 text-xs mt-1">Overdue</p>
              </div>
            </div>

            <div class="space-y-2.5">
              <p class="text-xs text-gray-900 font-medium tracking-widest">By priority</p>

              <div
                v-for="item in [
                  { label: 'High', value: taskStore.stats.byPriority.high, color: 'bg-red-600' },
                  { label: 'Medium', value: taskStore.stats.byPriority.medium, color: 'bg-yellow-600' },
                  { label: 'Low', value: taskStore.stats.byPriority.low, color: 'bg-gray-600' },
                ]"
                :key="item.label"
                class="flex items-center gap-3"
              >
                <span class="text-xs text-gray-900 w-12 flex-shrink-0">{{ item.label }}</span>
                <div class="flex-1 h-2 bg-green-200/90 rounded-full overflow-hidden">
                  <div
                    :class="['h-full rounded-full', item.color]"
                    :style="{
                      width: taskStore.stats.total > 0
                        ? `${Math.round((item.value / taskStore.stats.total) * 100)}%`
                        : '0%'
                    }"
                  ></div>
                </div>
                <span class="text-xs text-gray-900 w-4 text-right flex-shrink-0">{{ item.value }}</span>
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
