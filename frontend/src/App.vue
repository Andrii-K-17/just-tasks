<template>
  <RouterView />
  <footer class="flex justify-center shadow-md w-full border-t border-emerald-100 dark:border-slate-800 bg-emerald-50/30 dark:bg-slate-950 text-slate-900 dark:text-slate-400 text-center text-sm py-4">
    <span>&copy; {{ new Date().getFullYear() }} Just Tasks. All rights reserved.</span>
  </footer>
</template>
