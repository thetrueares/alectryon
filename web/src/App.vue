<template>
  <div class="min-h-screen bg-gray-100 flex flex-col items-center p-8">
    <nav class="mb-8 flex space-x-4">
      <router-link to="/" class="text-blue-600 hover:text-blue-800 font-medium">Home</router-link>
      <router-link to="/inputs" class="text-blue-600 hover:text-blue-800 font-medium">Inputs</router-link>
    </nav>
    <router-view></router-view>
  </div>
</template>

<style>
</style>
