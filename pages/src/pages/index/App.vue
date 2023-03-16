<template>
  <div v-if="logged">
    <UserInfo/>
  </div>
</template>

<script setup>
import {ref} from "vue";
import UserInfo from "./components/UserInfo.vue";

const logged = ref(false)

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}

if (!getCookie("us")) {
  location.href = "/login.html"
} else {
  logged.value = true
}

</script>

<style scoped>
</style>
