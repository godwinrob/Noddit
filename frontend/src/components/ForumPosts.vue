<template>

  <div class="forum-posts">
    
    <div v-for="post in postsArray" v-bind:key="post.postId">
      <post
        :postId="post.postId"
        :createdDate="post.createdDate"
        :username="post.username"
        :subnodditId="post.subnodditId"
        :subnodditName="post.subnodditName"
        :imageAddress="post.imageAddress"
        :title="post.title"
        :postScore="post.postScore"
        v-on:refresh-event="$emit('refresh-event')"
      />
    </div>

  </div>

</template>

<script>
import Post from "@/components/Post.vue";
//import API_Service from "@/service/API_Service.js";

export default {
  data() {
    return {};
  },
  props: {
    postsArray: Array
  },
  components: {
    Post
  },
  methods: {

  },
  created() {}
};
</script>

<style>
</style>