<template>
  <div class="home">
    <div class="shadow p-3 mb-5 bg-gray rounded">
    <div class="header-div"><h1>Today's popular posts</h1></div>
    <forum-posts :postsArray="recentPostsArray" v-on:refresh-event="getPopularPosts()"/>
    
    </div>
  </div>
</template>

<script>
import auth from '@/auth.js'
import API_Service from "@/service/API_Service.js";
import ForumPosts from "@/components/ForumPosts.vue"

export default {
  name: 'home',
  components: {
    ForumPosts
  },
  data() {
    return {
      user: auth.getUser(),
      recentPostsArray: []
    }
  },
  methods: {
    logout() {
      auth.logout();
      this.$router.push('/login');
    },
    getPopularPosts() {
      API_Service.getPopularPosts().then(
      parsedPosts => (this.recentPostsArray = parsedPosts)
    );
    }
  },
  created() {
    this.getPopularPosts();
  }
}
</script>

<style scoped>
.header-div {
  Margin: 5px;
  color: rgb(112, 112, 112);
}

</style>
