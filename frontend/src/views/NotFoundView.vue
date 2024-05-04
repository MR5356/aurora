<script setup lang="ts"></script>

<template>
  <div class="w-full h-full text-center flex justify-center flex-col items-center gap-10">
    <div />
    <img class="w-[20%] min-w-[200px] mt-[-100px]" src="@/assets/404-page-not-found.svg" alt="404" />
    <div class="text-sm" v-html="$t('404.message')"></div>
  </div>
</template>

<style scoped></style>
