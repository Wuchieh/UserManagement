<template>
  <div class="w-100 position-absolute text-center d-flex justify-content-center">
    <div id="liveAlertPlaceholder" class="col-12 col-md-7  col-xl-4 mx-0"></div>
  </div>
  <div
      class="container-fluid text-center vh-100 w-100 d-flex align-items-center justify-content-center position-sticky top-0">
    <form class="col-12 col-md-6 col-xl-3 p-3">
      <img alt="" class="mb-5 img-fluid" src="/img/newLogo.png" width="128">
      <ul class="nav nav-tabs">
        <li class="nav-item">
          <a :class="[page===0?'active':'']" aria-current="page" class="nav-link" href="#"
             @click.prevent="changePage(0)">登入</a>
        </li>
        <li class="nav-item">
          <a :class="[page===1?'active':'']" class="nav-link" href="#" @click.prevent="changePage(1)">註冊</a>
        </li>
        <li></li>
      </ul>
      <LoginPage v-if="page===0"/>
      <RegisterPage v-else @RegisterEd="RegisterEd"/>
      <p class="mt-5 mb-3 text-muted">© 2017–{{ year }}</p>
    </form>
  </div>
</template>

<script setup>
// const dt = new Date();
import LoginPage from "./components/LoginPage.vue";
import {ref} from "vue";
import RegisterPage from "./components/RegisterPage.vue";

const year = (new Date()).getFullYear()
const page = ref(0)

function changePage(i) {
  page.value = i
}

function RegisterEd() {
  page.value = 0
}
</script>

<style scoped>
#liveAlertPlaceholder {
  z-index: 4;
}
</style>

<style>
.alert {
  opacity: 0;
  animation: openIn 5s;
}

@keyframes openIn {
  0% {
    opacity: 0;
  }
  5% {
    opacity: 1;
  }
  95% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
</style>
