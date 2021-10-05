<template>
  <div class="root">
    <div class="login-container" ref="mainLoginContainer">
      <div class="login-box">
        <h2>{{ $t("login.login") }}</h2>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <input
              type="text"
              v-model="id"
              name="id"
              :placeholder="$t('login.id')"
              class="form-control"
              :class="{ 'is-invalid': submitted && !id }"
            />
            <div v-show="submitted && !id" class="form-error">
              {{ $t("login.id-required") }}
            </div>
          </div>
          <div class="form-group">
            <input
              type="password"
              v-model="password"
              name="password"
              :placeholder="$t('login.password')"
              class="form-control"
              :class="{ 'is-invalid': submitted && !password }"
            />
            <div v-show="submitted && !password" class="form-error">
              {{ $t("login.password-required") }}
            </div>
          </div>
          <div class="form-group" v-if="isRegister">
            <input
              type="text"
              v-model="name"
              name="name"
              :placeholder="$t('login.name')"
              class="form-control"
              :class="{ 'is-invalid': submitted && !name }"
            />
            <div v-show="submitted && !name && isRegister" class="form-error">
              {{ $t("login.name-required") }}
            </div>
          </div>
          <div class="form-group">
            <button
              class="btn btn-primary"
              :disabled="status.processing"
              v-if="!isRegister"
            >
              {{ $t("login.login") }}
            </button>
            <button
              class="btn btn-primary"
              :disabled="status.processing"
              v-if="isRegister"
            >
              {{ $t("login.register") }}
            </button>
            <img
              v-show="status.processing"
              class="processing"
              src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA=="
            />
            <a class="btn btn-link" @click="switchToLogin" v-if="isRegister">{{
              $t("login.login")
            }}</a>
            <a class="btn btn-link" @click="switchToRegister" v-if="!isRegister">{{
              $t("login.register")
            }}</a>
          </div>
          <div v-show="loginError" class="form-error">
            {{ loginError }}
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { SendGMCP } from "../gmcp";

export default {
  data() {
    return {
      id: "",
      password: "",
      name: "",
      submitted: false,
      isRegister: false,
      formError: "",
      status: {
        processing: false,
      },
    };
  },
  computed: {
    ...mapState(["isConnected", "gmcpOK", "isLogged", "loginError"]),
  },
  watch: {
    loginError: function (err) {
      if (err != "") {
        clearInterval(this.loginInterval);
        this.submitted = false;
        this.status.processing = false;
      }
    },
    isLogged: function (logged) {
      if (logged) {
        clearInterval(this.loginInterval);
        this.status.processing = false;
      }
    },
  },
  methods: {
    handleSubmit() {
      this.submitted = true;
      this.$store.state.loginError = "";

      let payload = {};
      if (!this.isRegister) {
        const { id, password } = this;
        if (!id || !password) {
          return;
        }
        payload = { id: id, password: password };
      } else {
        const { id, password, name } = this;
        if (!id || !password || !name) {
          return;
        }
        payload = {
          id: id,
          password: password,
          name: name,
        };
      }

      this.status.processing = true;
      if (!this.$store.state.isConnected) {
        this.$store.commit("CONNECT");
      }

      this.request(payload);
    },
    switchToLogin() {
      this.isRegister = false;
    },
    switchToRegister() {
      this.isRegister = true;
    },
    request(payload) {
      this.loginInterval = setInterval(
        function () {
          if (!this.$store.state.gmcpOK || this.$store.state.isLogged) {
            return;
          }
          if (!this.isRegister) {
            SendGMCP("Char.Login", payload);
          } else {
            SendGMCP("Char.Register", payload);
          }
        }.bind(this),
        1000
      );
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/styles/common.module";

.root {
  width: 100%;

  .login-container {
    text-align: center;

    .login-box {
      margin: 60px auto 0;
      width: 300px;
      text-align: left;
      background: $bg-color;

      h2 {
        text-align: center;
        font-size: 28px;
      }

      .form-error {
        font-size: 12px;
        color: salmon;
      }

      .form-group {
        margin: 7px 0;

        input {
          display: block;
          width: 100%;
          padding: 5px 8px;
          border: 0;
          color: $defaultTextColor;
          background-color: $bg-color-light3;
          font-family: $monoFont;
          font-weight: 500;
          font-size: 14px;
        }

        button {
          cursor: pointer;
          margin-right: 10px;
          padding: 4px 12px;
          color: #cacaca;
          font-weight: bold;
          background-color: $bg-color-light2;
          border-color: $bg-color-light2;
        }

        a {
          cursor: pointer;
          padding: 4px 0;
          font-size: 13px;
        }

        .processing {
          margin-right: 8px;
        }
      }
    }

    @media (max-width: 768px) {
      .login-box {
        width: 90%;
      }
    }
  }
}
</style>