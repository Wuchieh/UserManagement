<template>
  <hr>
  <div class="mx-2">
    <h1>新增用戶</h1>
    <div class="mb-3">
      <div class="mb-2">
        帳號: <input type="text" v-model="Username">
        密碼: <input type="text" v-model="Password">
      </div>
      <div class="">
        <input type="checkbox" v-model="Admin">是管理員
      </div>
    </div>
    <button class="btn btn-success" @click="sendCreateUser">送出</button>
  </div>
</template>

<script setup>
import {ref} from "vue";
import axios from "axios";

const Admin = ref(false)
const Username = ref("")
const Password = ref('')
const emits = defineEmits(['alertObj'])

function sendAlertObj(e) {
  emits('alertObj', e)
}

function sendCreateUser() {
  // console.log({
  //   Username: Username.value,
  //   Password: Password.value,
  //   Admin: Admin.value,
  // })

  if (Username.value.length < 1 || Password.value.length < 1) {
    sendAlertObj((Object({Message: "不可輸入空值", Type: "danger"})))
    return
  }

  axios({
    url: "/api/user",
    method: "POST",
    data: {
      Username: Username.value,
      Password: Password.value,
      Admin: Admin.value,
    }
  }).then((resp) => {
    sendAlertObj(Object({Message: resp.data["msg"], Type: "success"}))
  }).catch((err) => {
    if (err.response.status === 404) {
      sendAlertObj(Object({Message: "網路錯誤", Type: "danger"}))
    } else {
      sendAlertObj(Object({Message: err.response.data['msg'], Type: "danger"}))
    }
  })
}
</script>

<style scoped>

</style>