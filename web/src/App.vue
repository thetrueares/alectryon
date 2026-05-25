<template>
  <div class="flex min-h-screen bg-gray-100 text-gray-900">
    <!-- Sidebar -->
    <aside class="w-64 bg-black text-white shadow-xl flex-shrink-0">
      <div class="p-6">
        <h2 class="text-xl font-bold">Alectryon</h2>
      </div>
      <nav class="mt-4">
        <router-link 
          to="/" 
          class="block px-6 py-3 transition-colors duration-200 hover:bg-gray-800"
          active-class="bg-gray-900 border-r-4 border-blue-500"
        >
          Home
        </router-link>
        <router-link 
          to="/channels" 
          class="block px-6 py-3 transition-colors duration-200 hover:bg-gray-800"
          active-class="bg-gray-900 border-r-4 border-blue-500"
        >
          Channels
        </router-link>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-12">
      <router-view></router-view>
    </main>
  </div>
</template>

<style>
</style>
