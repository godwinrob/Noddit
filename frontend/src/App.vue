/**
  Noddit Messageboard
  Created by Rob Godwin, Craig Samad, Jay Minihan and Emma Knutson
  Capstone Project for Completion of Tech Elevator
  Version 1.0 Completed on 12/20/2019
 */
<template>
  <div id="app">
    // Header Nav Bar
    <b-navbar :sticky="true" id="nav" toggleable="lg" type="dark" variant="secondary">
      <b-navbar-brand href="/">
        <img src="@/assets/img/Logo-Text.png" alt="Image" class="logoText" />
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item href="/allsubnoddits">Browse Subnoddits</b-nav-item>
          <b-nav-item-dropdown>
            <template v-slot:button-content>
              <em>Active Subnoddits</em>
            </template>
            <b-dropdown-item
              v-for="sn in activeArray"
              :key="sn.subnodditName"
              :href="'/n/' + sn.subnodditName"
            >{{sn.subnodditName}}</b-dropdown-item>
          </b-nav-item-dropdown>
          <b-nav-item-dropdown v-show="getAuth() != null">
            <template v-slot:button-content>
              <em>Favorites</em>
            </template>
            <b-dropdown-item
              v-show="getAuth() != null && favoritesArray.length > 0"
              v-for="favorite in favoritesArray"
              :key="favorite.subnodditId"
              :href="'/n/' + favorite.subnodditName"
            >{{favorite.subnodditName}}</b-dropdown-item>
            <b-dropdown-item
              v-show="getAuth() != null && favoritesArray.length === 0"
              :href="'/allsubnoddits'"
            ><span class="no-favorites-text">No favorites yet...</span></b-dropdown-item>
          </b-nav-item-dropdown>
          <b-nav-item v-show="getAuth() != null" href="/newpost">New Post</b-nav-item>
          <b-nav-item v-show="getAuth() != null" href="/newSubnoddit">New Subnoddit</b-nav-item>
        </b-navbar-nav>

          <b-navbar-nav class="ml-auto">
            <b-nav-form>
              <b-form-input class="mr-sm-2" placeholder="Search for Subnoddits" v-model="searchTerm"></b-form-input>
              <button class="btn btn-dark" @click.prevent="searchSubnoddits">Search</button>
            </b-nav-form>

            <b-nav-item href="/login" v-show="getAuth() === null">Login</b-nav-item>
            <b-nav-item-dropdown right v-show="getAuth() != null">
              
              <template v-slot:button-content>
                <em>{{getUsername()}}</em>
              </template>
              <b-dropdown-item href="/profile">Profile</b-dropdown-item>
              <b-dropdown-item @click="logout">Sign Out</b-dropdown-item>
            </b-nav-item-dropdown>
          </b-navbar-nav>
       
      </b-collapse>
    </b-navbar>

    // Page Content
    <div class="app-body">
      <router-view />
    </div>

    // Footer 
    <div class="footer">
      <page-footer />
    </div>
  </div>
</template>

<script>
import auth from "@/auth.js";
import PageFooter from "@/components/PageFooter.vue";
import API_Service from "@/service/API_Service.js";

export default {
  components: {
    PageFooter
  },
  data() {
    return {
      user: auth.getUser(),
      favoritesArray: [],
      activeArray: [],
      searchTerm: ''
    };
  },
  created() {
    if (this.getUsername() != null) {
      API_Service.getFavoritesForUser(this.getUsername()).then(
        parsedFavorites => (this.favoritesArray = parsedFavorites)
      );
    }
    API_Service.getActiveSubnoddits().then(
      parsedSubnoddits => (this.activeArray = parsedSubnoddits)
    );
  },
  methods: {
    logout() {
      auth.logout();
      this.$router.push("/login");
      this.refreshPage();
    },
    getAuth() {
      return auth.getUser();
    },
    getUsername() {
      if (auth.getUser() == null) {
        return null;
      } else {
        return auth.getUser().sub;
      }
    },
    refreshPage() {
      window.location.reload();
    },
    searchSubnoddits() {
      this.$router.push("/search/" + this.searchTerm)
      this.refreshPage();
    }
  }
};
</script>

<style>
#app {
  background-color: lightgray;
}

.logoText {
  width: 6rem;
}
.no-favorites-text {
  color: gray;
  font-style: italic;
}

.footer {
  text-align: center;
  bottom: 0;
}
.app-body {
  min-height: 500px;
}
</style>