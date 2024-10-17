<template>
  <NuxtLayout>
    <NuxtPage/>
    <Toaster richColors position="top-center"/>
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
import '~/assets/simple-markdown.scss'
</script>

<style>
div.grecaptcha-badge{
  display: none !important;
}
</style>