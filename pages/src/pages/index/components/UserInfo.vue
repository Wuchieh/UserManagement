<template>
  <!-- Loading -->
  <Loading v-if="LoadStatus === 'Loading'"/>

  <!-- Error Page -->
  <Error v-else-if="LoadStatus==='Error'"/>

  <!-- Main Page-->
  <Main v-else-if="LoadStatus === 'Success'" :Username="Username" :admin="admin"/>
</template>

<script setup>
import axios from "axios";
import {ref} from "vue";
import Loading from "./Loading.vue";
import Error from "./Error.vue";
import Main from "./Main.vue";

const Username = ref("")
const LoadStatus = ref("Loading")
const admin = ref(false)

// 發起網路請求 請求用戶資訊 並且回傳給父模板
axios({
  url: "/api/user/userInfo"
}).then((resp) => {
  Username.value = String(resp.data["Username"]) + " "
  admin.value = Boolean(resp.data["Admin"])
  setTimeout(() => {
    LoadStatus.value = 'Success'
  }, 300)
}).catch((err) => {
  console.log(err)
  LoadStatus.value = 'Error'
  if (err.response.data['status'] === false) {
    location.href = "/login.html";
  }
})
</script>

<style scoped>
@font-face {
  font-family: 'Helvetica';
  src: url('/src/assets/fonts/Nunito-SemiBold.ttf')
}

/*#myCanvas {*/
/*  height: 50px;*/
/*  width: 150px;*/
/*  background: #afafaf;*/
/*}*/

/*#canvasText {*/
/*  font-family: 'Helvetica',serif ;*/
/*  margin-top: 13px;*/
/*  font-size: 25px;*/
/*}*/
</style>