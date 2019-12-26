<template>


  <div id="login-form">
    <div class="shadow-lg p-3 mb-5 bg-white rounded">
    <div id="login" class="text-center">
      
      <b-card>
      <form class="form-signin" @submit.prevent="login">
        <h1 class="h3 mb-3 font-weight-normal">Please Sign In</h1>
        <div
          class="alert alert-danger"
          role="alert"
          v-if="invalidCredentials"
        >Invalid username and password!</div>
        <div
          class="alert alert-success"
          role="alert"
          v-if="this.$route.query.registration"
        >Thank you for registering, please sign in.</div>
        <label for="username" class="sr-only">Username</label>
        <input
          type="text"
          id="username"
          class="form-control"
          placeholder="Username"
          v-model="user.username"
          required
          autofocus
        />
        <label for="password" class="sr-only">Password</label>
        <input
          type="password"
          id="password"
          class="form-control"
          placeholder="Password"
          v-model="user.password"
          required
        />
        <router-link :to="{ name: 'register' }">Need an account?</router-link>
        <br />
        <gb-button class="btn btn-lg btn-primary btn-block" type="submit" color="black">Sign in</gb-button>
      </form>
      </b-card>
      </div>
    </div>
  </div>

</template>

<script>
import auth from "../auth";

export default {
  name: "login",
  components: {},
  data() {
    return {
      user: {
        username: "",
        password: ""
      },
      invalidCredentials: false
    };
  },
  methods: {
    login() {
      fetch(`${process.env.VUE_APP_REMOTE_API}/login`, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json"
        },
        body: JSON.stringify(this.user)
      })
        .then(response => {
          if (response.ok) {
            return response.text();
          } else {
            this.invalidCredentials = true;
          }
        })
        .then(token => {
          if (token != undefined) {
            if (token.includes('"')) {
              token = token.replace(/"/g, "");
            }
            auth.saveToken(token);
            this.$router.push("/");
            this.refreshPage();
          }
        })
        .catch(err => console.error(err));
    },
    refreshPage() {
      window.location.reload();
    }
  }
};
</script>

<style scoped>
#login-form {
  margin-top: 100px;
  display: flex;
  justify-content: center;
}
#login {
  border-radius: 5px;
  padding: 10px;
  background-color: rgb(245, 245, 245);
  border: 1px solid rgb(197, 197, 197);
  max-width: 400px;
}
</style>
