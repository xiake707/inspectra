<script setup>
import { reactive, ref } from "vue";
import { hostRules } from "../utils/hostValidation.js";

const visible = defineModel({ type: Boolean, default: false });
const emit = defineEmits(["submit"]);
const formRef = ref();
const form = reactive({
  Host: "",
  User: "root",
  Port: 22,
  Password: "",
  KeyFile: "",
});

function reset() {
  Object.assign(form, {
    Host: "",
    User: "root",
    Port: 22,
    Password: "",
    KeyFile: "",
  });
  formRef.value?.clearValidate();
}

async function submit() {
  try {
    await formRef.value.validate();
    emit("submit", form);
  } catch {
    /* Element Plus renders field errors. */
  }
}
</script>
<template>
  <el-dialog
    v-model="visible"
    width="570px"
    title="添加主机"
    align-center
    destroy-on-close
    @closed="reset"
  >
    <div class="dialog-intro">
      填写目标主机的 SSH 连接信息，保存前会进行格式校验。
    </div>
    <el-form
      ref="formRef"
      :model="form"
      :rules="hostRules"
      label-position="top"
      class="host-form"
    >
      <el-row :gutter="16">
        <el-col :span="24">
          <el-form-item label="主机地址" prop="Host">
            <el-input
              v-model="form.Host"
              placeholder="例如 192.168.122.18"
              maxlength="15"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="SSH 用户" prop="User">
            <el-input v-model="form.User" placeholder="root" maxlength="20" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="SSH 端口" prop="Port">
            <el-input-number
              v-model="form.Port"
              :min="1"
              :max="65535"
              controls-position="right"
              class="full-width"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="登录密码" prop="Password">
            <el-input
              v-model="form.Password"
              type="password"
              show-password
              placeholder="可选"
              maxlength="20"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="私钥路径" prop="KeyFile">
            <el-input
              v-model="form.KeyFile"
              placeholder="/root/.ssh/id_ed25519"
              maxlength="30"
            />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submit">生成请求</el-button>
    </template>
  </el-dialog>
</template>
