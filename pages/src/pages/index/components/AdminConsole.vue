<template>
  <Alert :msg="alertMsg"/>

  <AdminCreateUser @alertObj="sendAlert"/>
  <div v-show="editor">
    <AdminEditor :user="selectUser" @closeEditor="editorClose" @reloadData="reloadTable"/>
  </div>

  <!-- 確認刪除 -->
  <div id="deleteUser" style="display: none">
    <div class="bg-dark vh-100 vw-100 top-0 start-0 opacity-25 position-absolute" @click="cancelDeleteUser"
         style="z-index: 1002"></div>
    <div class="position-absolute d-flex justify-content-center align-items-center w-100 h-100">
      <div class="p-2 border-dark border rounded-2 bg-white" style="z-index: 1004">
        <div class="">確定要刪除此筆資料嗎?</div>
        <div class="mb-3">Username: {{ selectUser['Username'] }}</div>

        <button class="btn btn-danger me-3" @click="ConfirmDeleteUser">確定</button>
        <button class="btn btn-secondary" @click="cancelDeleteUser">取消</button>
      </div>
    </div>
  </div>

  <!-- 右鍵選單 -->
  <div id="rightMenu" class="border border-dark rounded-2">
    <div class="menu m-1">
      <button class="btn btn-outline-secondary col-12 mb-1" @click="enableEditor">修改</button>
      <button class="btn btn-outline-danger col-12" @click="deleteUser">刪除</button>
    </div>
  </div>

  <hr>
  <!-- 控制台主頁面 -->
  <div class="d-flex mb-1">
    <h1>管理員控制台</h1>

    <div class="ms-3">
      <span>每頁幾筆資料: </span>
      <select v-model="limit">
        <option v-for="i in [10,20,50]" :selected="limit===i" :value="i">{{ i }}</option>
      </select>
    </div>

    <button class="btn btn-primary position-absolute end-0 me-3 align-baseline d-none d-md-block" @click="reloadTable">
      <svg class="reload bi bi-arrow-repeat" fill="currentColor" height="20" viewBox="0 0 16 16"
           width="20" xmlns="http://www.w3.org/2000/svg">
        <path
            d="M11.534 7h3.932a.25.25 0 0 1 .192.41l-1.966 2.36a.25.25 0 0 1-.384 0l-1.966-2.36a.25.25 0 0 1 .192-.41zm-11 2h3.932a.25.25 0 0 0 .192-.41L2.692 6.23a.25.25 0 0 0-.384 0L.342 8.59A.25.25 0 0 0 .534 9z"/>
        <path
            d="M8 3c-1.552 0-2.94.707-3.857 1.818a.5.5 0 1 1-.771-.636A6.002 6.002 0 0 1 13.917 7H12.9A5.002 5.002 0 0 0 8 3zM3.1 9a5.002 5.002 0 0 0 8.757 2.182.5.5 0 1 1 .771.636A6.002 6.002 0 0 1 2.083 9H3.1z"
            fill-rule="evenodd"/>
      </svg>
      刷新
    </button>
  </div>

  <!-- 資料展示區 -->
  <div v-if="debug">
    <div class="table-responsive ms-2 mb-2">
      <table class="table table-striped">
        <thead>
        <tr>
          <th scope="col">ID</th>
          <th scope="col">Username</th>
          <th scope="col">Password</th>
          <th scope="col">CreateTime</th>
          <th scope="col">UpdateTime</th>
          <th scope="col">Admin</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="i in UserList" @click.right.prevent="clickRight($event,i)">
          <th class="data" scope="row">{{ i['ID'] }}</th>
          <td class="data">{{ i['Username'] }}</td>
          <td :title="i['Password']" class="data">{{ truncate(i['Password'], 20) }}</td>
          <td class="data">{{ dateFormat(i['CreateTime']) }}</td>
          <td class="data">{{ dateFormat(i['UpdateTime']) }}</td>
          <td class="data">{{ i['Admin'] ? '是' : '否' }}</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- 頁面切換 -->
    <nav aria-label="Page navigation example" class="d-flex">
      <ul class="pagination">
        <input :value="page" class="m-2 mt-0" style="width: 5em;height: 100%" type="number"
               @change="inputPage" @keydown.enter="inputPage">

        <li class="page-item">
          <button class="page-link" @click="changeUserListPage('first')">
            <span aria-hidden="true">&laquo;</span>
          </button>
        </li>
        <li class="page-item">
          <button class="page-link" @click="changeUserListPage('previous')">
            <span aria-hidden="true">&lsaquo;</span>
          </button>
        </li>

        <li v-for="i in createPageButton()" class="page-item">
          <button :class="[i===page?'active':'']"
                  :disabled="i===page"
                  class="page-link"
                  @click="changeUserListPage(i)">
            {{ i }}
          </button>
        </li>

        <li class="page-item" @click="changeUserListPage('next')">
          <button class="page-link">
            <span aria-hidden="true">&rsaquo;</span>
          </button>
        </li>
        <li class="page-item" @click="changeUserListPage('last')">
          <button class="page-link">
            <span aria-hidden="true">&raquo;</span>
          </button>
        </li>

      </ul>
    </nav>
  </div>
</template>

<script setup>
import axios from "axios";
import {onMounted, ref, watch} from "vue";
import AdminEditor from "./AdminEditor.vue";
import Alert from "./Alert.vue";
import AdminCreateUser from "./AdminCreateUser.vue";

const debug = true
const UserList = ref([Object])
const page = ref(1)
const maxPage = ref(0)
const limit = ref(10)
const selectUser = ref("")
const editor = ref(false)
const alertMsg = ref({Message: "", Type: ""})

function deleteUser() {
  document.getElementById("deleteUser").style.display = "block";
}

function cancelDeleteUser() {
  document.getElementById("deleteUser").style.display = "none";
}

function sendAlert(e) {
  alertMsg.value = e
  if(e.Type==="success"){
    reloadTable()
  }
}

function ConfirmDeleteUser() {
  axios({
    url: "/api/user",
    method: "DELETE",
    data: selectUser.value
  }).then(() => {
    reloadTable()
    cancelDeleteUser()

  }).catch(err => {
    if (err.response.status === 404) {
      sendAlert({
        Message: "網路錯誤",
        Type: "danger"
      })
    } else {
      sendAlert({
        Message: err.response.data["msg"],
        Type: "danger"
      })
    }
  })
}

function dateFormat(date) {
  let dateA = new Date(String(date))
  if (dateA.getDate()) {
    return dateA.getFullYear() + "-" + (dateA.getMonth() + 1) + "-" + dateA.getDate() +
        " " + dateA.getHours() + ":" + dateA.getMinutes() + ":" + dateA.getSeconds()
  }
}

function truncate(value, length) {
  if (!value) return ""
  value = value.toString()
  if (value.length < length) {
    return value
  }
  return value.substring(0, length) + '…'
}

function getUserList() {
  axios({
    url: "/api/users/" + (page.value - 1),
    params: {
      limit: limit.value
    }
  }).then((resp) => {
    UserList.value = resp.data["userList"]
    maxPage.value = Math.ceil((resp.data["countDocuments"] / limit.value))
  })
}

function changeUserListPage(pageNum) {
  switch (pageNum) {
    case "first":
      page.value = 1
      break
    case "previous":
      if (page.value <= 1) {
        return
      }
      page.value--
      break
    case "last":
      page.value = maxPage.value
      break
    case "next":
      if (page.value === maxPage.value) {
        return
      }
      page.value++
      break
    default:
      page.value = pageNum
      break
  }
}

function inputPage(e) {
  page.value = Math.trunc(e.target.value)
}

function createPageButton() {
  const buttons = [];
  const numButtons = Math.min(maxPage.value, 5);
  const startPage = Math.max(1, Math.min(page.value - 2, maxPage.value - 4));

  for (let i = 0; i < numButtons; i++) {
    buttons.push(startPage + i);
  }

  return buttons;
}


function reloadTable() {
  axios({
    url: "/api/user/getUserCount"
  }).then(() => {
    getUserList()
  })
}

function clickRight(e, i) {
  let rightMenu = document.getElementById("rightMenu")
  rightMenu.style.display = "block"
  rightMenu.style.top = e.clientY + 'px'
  rightMenu.style.left = e.clientX + 'px'
  selectUser.value = i
}

function enableEditor() {
  editor.value = true
}

function editorClose() {
  editor.value = false
}

getUserList()
watch([page, limit], () => {
  getUserList()
})
onMounted(() => {
  let rightMenu = document.getElementById("rightMenu")
  document.onclick = (e) => {
    if (e['target'] === rightMenu) {
      return
    }
    rightMenu.style.display = "none"
  }
})

</script>

<style scoped>
.reload {
  transform: rotate(-36deg);
  -webkit-transform: rotate(-36deg);
}

#rightMenu {
  position: absolute;
  display: none;
  width: 5em;
  background: rgb(255, 255, 255);
  z-index: 1000;
}

#deleteUser {
  position: absolute;
  display: none;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  z-index: 100;
}
</style>