<template>
  <alert style="z-index: 1005" :msg="alertMsg"/>
  <div class="editorBg" @click="close">
  </div>

  <div class="align-items-center justify-content-center d-flex editor border border-dark bg-light">
    <div class="p-2">
      <div class="d-flex flex-row-reverse mb-3">
        <button class="btn btn-danger rounded-5" @click="close">X</button>
      </div>
      <div>
        <div class="input-group mb-3">
          <span class="input-group-text">Username</span>
          <input type="text" class="form-control" disabled v-model="editUsername">
        </div>
        <div class="input-group">
          <span class="input-group-text">Password</span>
          <input type="text" class="form-control" v-model="editPassword">
        </div>
        <div class="mb-3">
          <span class="text-danger">密碼會自動加密 請輸入明文 </span><br>
          <span class="text-danger">長度小於8則不修改</span>
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">Admin</span>
          <select class="form-select form-select-sm" v-model="editAdmin">
            <option :value="true">是</option>
            <option :value="false">否</option>
          </select>
        </div>
      </div>

      <button class="btn btn-success fs-4" @click="wantEdit">修改</button>
    </div>
  </div>

  <div class="checkNotify text-center" v-show="editCheck">
    <div class="my-2">
      <h1>確定要進行修改?</h1>
      <div>
        <button class="btn btn-primary me-3 fs-5" @click="ConfirmEdit" :disabled="axiosStatus">
          <div class="spinner-border" role="status" style="height: 1rem;width: 1rem" v-show="axiosStatus">
            <span class="visually-hidden">Loading...</span>
          </div>
          確定
        </button>
        <button class="btn btn-secondary fs-5" @click="CancelEdit">取消</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, watch} from "vue";
import axios from "axios";
import Alert from "./Alert.vue";

const props = defineProps(["user"])
const emit = defineEmits(["closeEditor","reloadData"])

const editCheck = ref(false)
const editUsername = ref("")
const editPassword = ref("")
const editAdmin = ref(false)
const editUserId = ref("")
const axiosStatus = ref(false)
const alertMsg = ref({Message:"",Type:""})

function wantEdit() {
  editCheck.value = true
}

function CancelEdit() {
  editCheck.value = false
}

function close() {
  emit('closeEditor')
  editCheck.value = false
}

function reload(){
  emit("reloadData")
}

function axiosStatusFalse() {
  axiosStatus.value = false
}

function ConfirmEdit() {
  // console.log(
  //     {
  //       Username:editUsername.value,
  //       Password:editPassword.value,
  //       Admin:editAdmin.value,
  //     }
  // )
  axiosStatus.value = true
  axios({
    url: "/api/user",
    method: "PUT",
    data: {
      ID: editUserId.value,
      Password: editPassword.value,
      Admin: editAdmin.value,
    },
  }).then((resp) => {
    axiosStatusFalse()
    if (resp['data']['status']) {
      alertMsg.value = {Message:resp['data']['msg'],Type:"success"}
      setTimeout(() => {
        alertMsg.value = {Message:"",Type:""}
        close()
        reload()
      }, 1000)

    }
  }).catch((err) => {
    axiosStatusFalse()
    console.log(err)
    if (err.response.status !== 404) {
      alertMsg.value = {Message:err.response['data']['msg'],Type:"danger"}
      setTimeout(() => {
        alertMsg.value = {Message:"",Type:""}
      }, 100)

    } else {
      alertMsg.value = {Message:"網路錯誤",Type:"danger"}
      setTimeout(() => {
        alertMsg.value = {Message:"",Type:""}
      }, 100)

    }
  })
}

watch(props, (n) => {
  let user = n['user'];
  editUsername.value = user["Username"]
  editAdmin.value = user["Admin"]
  editUserId.value = user["ID"]
})
</script>

<style scoped>
.editorBg {
  position: absolute;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  background: rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.editor {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1002;
  border-radius: 10px;
}

.checkNotify {
  position: absolute;
  background: white;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100vw;
  z-index: 1004;
  animation: ease-in-out .25s;
}

@keyframes ease-in-out {
  0% {
    top: 75%
  }
  100% {
    top: 50%
  }
}
</style>