<template>
  <div class="profile-grid">
    <profile-form class="accountCard" />

   <div class="myPosts">
<div class="shadow-lg p-3 mb-5 bg-gray rounded">
        <div class="header-div"><h1>My Posts</h1></div>
          <forum-posts :postsArray="recentPostsArray" v-on:refresh-event="getPopularPosts()"/>
         


          <!--<div>
        <ul class="nav nav-tabs">
          <li class="nav-item">
            <a class="nav-link" href="#">My Posts</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Liked</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Favorite Forums</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#" >Following</a>
          </li>
        </ul>
     <forum-posts/>-->
    </div>
  </div>
  </div>
</template>

<script>
import auth from '@/auth.js'
import API_Service from "@/service/API_Service.js";
import ForumPosts from "@/components/ForumPosts.vue"
import ProfileForm from "@/components/ProfileForm.vue";

export default {
  name: 'profile',
  components: {
    ProfileForm,
    ForumPosts
  },
data() {
    return {
      user: auth.getUser(),
      recentPostsArray: []
    }
  },
  methods:{
      getPopularPosts() {
      console.log('here');
      API_Service.getPostsForUser(auth.getUser().sub).then(
      parsedPosts => (this.recentPostsArray = parsedPosts)
    );
    }
  },
  created() {
    this.getPopularPosts();
  }
}

</script>

<style>
.profile-grid {
  display: grid;
  grid-template-columns: 0.5fr 5fr 0.25fr 2fr 0.5fr;
  grid-template-areas: " . myPosts . accountCard . ";
}
.myPosts {
  grid-area: myPosts;
  margin-top: 50px;
}
.header-div {
  color:rgb(112, 112, 112);
}
.accountCard {
  grid-area: accountCard;
  margin-top: 3rem;

}
</style>