<template>
  <div class="m-2">
    <div class="d-flex">
      <h1>{{ Username }}Hello</h1>
      <input class="btn btn-danger ms-3" type="button" value="登出" @click="logout">
    </div>
  </div>
  <AdminConsole v-if="props['admin']"/>
</template>

<script setup>
import axios from "axios";
import AdminConsole from "./AdminConsole.vue";

const props = defineProps(["Username", "admin"])

const logout = () => {
  axios({
    url: "/api/user/logout",
    method: "POST"
  })
  location.href = location.pathname
}

if (props["admin"]) {
  console.log("AdminLogin")
}

</script>

<style scoped>

</style>