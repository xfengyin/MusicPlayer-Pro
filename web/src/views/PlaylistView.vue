<template>
  <div class="playlist-view">
    <h1>Playlist</h1>
    <ul>
      <li v-for="song in songs" :key="song.id">{{ song.title }}</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
const songs = [
  { id: 1, title: 'Sample Song 1' },
  { id: 2, title: 'Sample Song 2' }
]
</script>

<style scoped>
.playlist-view {
  padding: 20px;
}
</style>
