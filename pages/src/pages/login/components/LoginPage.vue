<template>

  <div class="form-floating">
    <input id="floatingInput" v-model="Username" class="form-control" placeholder="Username" type="text">
    <label for="floatingInput">帳號</label>
  </div>

  <div class="form-floating">
    <input id="floatingPassword" v-model="Password" autocomplete="on" class="form-control" placeholder="Password"
           type="password">
    <label for="floatingPassword">密碼</label>
  </div>

  <div class="checkbox mb-3">
    <label>
      <input v-model="Checked" type="checkbox" value="remember-me"> 記住我 7天
    </label>
  </div>

  <button class="w-100 btn btn-lg btn-primary" type="submit" @click.prevent="login">
    <div v-show="apiStatus" class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    登入
  </button>
</template>

<script setup>
import {ref} from "vue";
import axios from "axios";

const Username = ref("")
const Password = ref("")
const apiStatus = ref(false)
const Checked = ref(false)
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
  }
}

function login(e) {
  if (!(Password.value.length || Username.value.length)) {
    alertShow("請輸入帳號和密碼", "danger")
    return
  }
  if (apiStatus.value) {
    return;
  }
  e.target.classList.add('disabled')
  apiStatus.value = true

  axios({
    method: "POST",
    url: "/api/user/login",
    data: {
      Username: Username.value,
      Password: Password.value,
      Checked: Checked.value
    }
  }).then(
      (response) => {
        alert(response.data, e)
        setTimeout(() => {
          location.href = "/"
        }, 1000)
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

.form-floating input[type="password"] {
  border-top-right-radius: 0;
  border-top-left-radius: 0;
  margin-bottom: 10px;
}

.form-floating:focus-within {
  z-index: 2;
}

.spinner-border {
  --bs-spinner-width: 1em;
  --bs-spinner-height: 1em;
}
</style>