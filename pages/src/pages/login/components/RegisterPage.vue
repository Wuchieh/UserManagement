<template>

  <div class="form-floating">
    <input id="floatingInput" v-model="Username" class="form-control" placeholder="Username" type="text">
    <label for="floatingInput">帳號</label>
  </div>

  <div class="form-floating">
    <input id="floatingPassword" v-model="Password" autocomplete="on" class="form-control" name="Password"
           placeholder="Password" type="password">
    <label for="floatingPassword">密碼</label>
  </div>

  <div class="form-floating">
    <input id="floatingConfirmPassword" v-model="ConfirmPassword" autocomplete="on" class="form-control"
           name="ConfirmPassword"
           placeholder="Password" type="password">
    <label for="floatingConfirmPassword">確認密碼</label>
  </div>

  <button class="w-100 btn btn-lg btn-primary" @click.prevent="register">
    <div v-show="apiStatus" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    註冊
  </button>
</template>

<script setup>
import {ref} from "vue";
import axios from "axios";

const Username = ref("")
const Password = ref("")
const ConfirmPassword = ref("")
const apiStatus = ref(false)


const emits = defineEmits(["RegisterEd"])

// type:{danger,success}
const alertShow = (message, type) => {
  const wrapper = document.createElement('div')
  wrapper.innerHTML = [
    `<div class="alert alert-${type} alert-dismissible mb-2" role="alert">`,
    `   <div>${message}</div>`,
    '   <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>',
    '</div>'
  ].join('')
  const alertPlaceholder = document.getElementById('liveAlertPlaceholder')
  alertPlaceholder.append(wrapper)
  setTimeout(() => {
    alertPlaceholder.removeChild(wrapper)
  }, 4999)
}

function alert(data, e) {
  let type = data['status'] ? `success` : `danger`
  alertShow(data['msg'], type)

  if (type !== `success`) {
    e.target.classList.remove('disabled')
    apiStatus.value = false
  } else {
    emits("RegisterEd")
  }
}

function register(e) {
  if (Password.value !== ConfirmPassword.value) {
    alertShow("兩次密碼不同", "danger")
    return
  }
  if (apiStatus.value) {
    return;
  }
  e.target.classList.add('disabled')
  apiStatus.value = true

  axios({
    method: "POST",
    url: "/api/user/register",
    data: {
      Username: Username.value,
      Password: Password.value
    }
  }).then(
      (response) => {
        alert(response.data, e)
      }
  ).catch(
      (error) => {
        if (error.response.status === 404) {
          alert({msg: "網路發生問題請稍後再試"}, e)
        } else {
          alert(error.response.data, e)
        }
      }
  )
}


</script>

<style scoped>
.form-floating input[type="text"] {
  border-radius: 0;
  margin-bottom: -1px;
}

.form-floating input[name="Password"] {
  border-radius: 0;
  margin-bottom: -1px;
}

.form-floating input[name="ConfirmPassword"] {
  border-top-right-radius: 0;
  border-top-left-radius: 0;
  margin-bottom: 20px;
}

.form-floating:focus-within {
  z-index: 2;
}

.spinner-border {
  --bs-spinner-width: 1em;
  --bs-spinner-height: 1em;
}
</style>