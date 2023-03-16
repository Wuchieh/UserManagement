<template>
  <div class="w-100 position-absolute text-center d-flex justify-content-center">
    <div id="liveAlertPlaceholder" class="col-12 col-md-7  col-xl-4 mx-0"></div>
  </div>
  <div
      class="container-fluid text-center vh-100 w-100 d-flex align-items-center justify-content-center position-sticky top-0">
    <div class="">
      <img alt="" class="mb-1 img-fluid" src="/img/newLogo.png" width="128">
      <hr>
      <h1 class="">發生例外錯誤！！！</h1>
      <h1>請稍後刷新頁面</h1>
      <button class="btn btn-primary p-2" @click="reload">
        <svg class="reload bi bi-arrow-repeat" fill="currentColor" height="30" viewBox="0 0 16 16"
             width="30" xmlns="http://www.w3.org/2000/svg">
          <path
              d="M11.534 7h3.932a.25.25 0 0 1 .192.41l-1.966 2.36a.25.25 0 0 1-.384 0l-1.966-2.36a.25.25 0 0 1 .192-.41zm-11 2h3.932a.25.25 0 0 0 .192-.41L2.692 6.23a.25.25 0 0 0-.384 0L.342 8.59A.25.25 0 0 0 .534 9z"/>
          <path d="M8 3c-1.552 0-2.94.707-3.857 1.818a.5.5 0 1 1-.771-.636A6.002 6.002 0 0 1 13.917 7H12.9A5.002 5.002 0 0 0 8 3zM3.1 9a5.002 5.002 0 0 0 8.757 2.182.5.5 0 1 1 .771.636A6.002 6.002 0 0 1 2.083 9H3.1z"
                fill-rule="evenodd"/>
        </svg>
      </button>
      <hr>
      <p class="mt-2 text-muted">© 2017–{{ year }}</p>
    </div>
  </div>
</template>

<script setup>
const year = (new Date()).getFullYear()
const reload = () => {
  location.href = location.pathname
}
</script>

<style scoped>
.reload {
  transform: rotate(-36deg);
  -webkit-transform: rotate(-36deg);
}
</style>