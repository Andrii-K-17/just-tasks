<script setup lang="ts">
import { useRouter } from 'vue-router'
import { 
  CheckCircleIcon,
  ChartPieIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen bg-emerald-50/30 text-black flex flex-col font-sans">
    
    <div class="sticky h-14 top-0 z-50 shadow-xs w-full border-b border-emerald-100 bg-[#f9fefb]">
      <header class="flex items-center justify-between max-w-5xl mx-auto px-4 h-full">
        <div>
          <h1 class="text-xl font-bold tracking-tight inline-flex items-center gap-1 cursor-pointer" @click="router.push('/')">
            Just <span class="text-emerald-500">Tasks</span>
          </h1>
        </div>

        <div class="flex items-center gap-2 sm:gap-4 text-sm font-medium">
          <button 
            @click="router.push('/login')" 
            class="px-4 py-2 rounded-xl border border-emerald-400 text-slate-700 hover:bg-emerald-100 transition-colors cursor-pointer"
          >
            Log in
          </button>
          <button 
            @click="router.push('/login?mode=register')" 
            class="px-4 py-2 rounded-xl bg-emerald-600 text-white hover:bg-emerald-700 transition-colors shadow-sm cursor-pointer"
          >
            Get Started
          </button>
        </div>
      </header>
    </div>

    <main class="flex-grow flex flex-col items-center justify-center px-4 pt-13 pb-16 text-center">
      <h2 class="text-5xl sm:text-6xl md:text-7xl font-extrabold text-slate-900 mb-6 max-w-4xl">
        Organize your tasks,
        <span class="text-emerald-500 text-transparent bg-clip-text bg-gradient-to-r from-emerald-500 to-emerald-700">
          without the chaos.
        </span>
      </h2>
      
      <p class="text-lg sm:text-xl text-gray-600 max-w-2xl mb-10 leading-relaxed">
        A minimalist space for your daily tasks.
      </p>

      <button 
        @click="router.push('/login?mode=register')" 
        class="px-8 py-3.5 rounded-full bg-emerald-600 text-white text-lg font-medium hover:bg-emerald-700 hover:shadow-lg hover:shadow-emerald-600/20 transform hover:-translate-y-0.5 transition-all cursor-pointer"
      >
        Start doing for free
      </button>

    </main>

    <section class="bg-[#f9fefb] border-t border-emerald-100 py-10 px-4">
      <div class="max-w-4xl mx-auto grid sm:grid-cols-2 md:grid-cols-2 gap-10 sm:gap-12">
        
        <div class="space-y-4">
          <div class="w-12 h-12 rounded-2xl bg-emerald-100 text-emerald-600 flex items-center justify-center">
            <CheckCircleIcon class="w-6 h-6" />
          </div>
          <h3 class="text-xl font-bold text-slate-900">Simple by design</h3>
          <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
            A clean interface that respects your attention and time.
          </p>
        </div>

        <div class="space-y-4">
          <div class="w-12 h-12 rounded-2xl bg-emerald-100 text-emerald-600 flex items-center justify-center">
            <ChartPieIcon class="w-6 h-6" />
          </div>
          <h3 class="text-xl font-bold text-slate-900">Clear progress</h3>
          <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
            Understand your productivity with minimal stats.
          </p>
        </div>

      </div>
    </section>
  </div>
</template>

<style scoped>
.shadow-xs {
  box-shadow: 0 1.5px 1px rgba(130, 150, 133, 0.05);
}
</style>
