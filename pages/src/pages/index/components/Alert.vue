<template>
  <div class="w-100 position-absolute text-center d-flex justify-content-center top-0">
    <div id="liveAlertPlaceholder" class="col-12 col-md-7  col-xl-4 mx-0">
    </div>
  </div>
</template>

<script setup>
import {watch} from "vue";

const prop = defineProps(["msg"])
watch(prop, (v) => {
  if (v.msg["Message"] !== '') {
    alertShow(v.msg["Message"],v.msg["Type"])
  }
})

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
  }, 999)
}
</script>

<style>

.alert {
  opacity: 0;
  animation: openIn 1s;
  z-index: 2000;
}

@keyframes openIn {
  0% {
    opacity: 0;
    top: 10px;
  }
  5% {
    opacity: 1;
  }
  10%{
    top: 0;
  }
  95% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
</style>