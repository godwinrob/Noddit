<template>
  <div>
    <post
      :postId="parentPost.postId"
      :createdDate="parentPost.createdDate"
      :username="parentPost.username"
      :subnodditId="parentPost.subnodditId"
      :subnodditName="parentPost.subnodditName"
      :imageAddress="parentPost.imageAddress"
      :title="parentPost.title"
      :postScore="parentPost.postScore"
      :body="parentPost.body"
      :commentPage="true"
      v-on:refresh-event="populateRepliesArray()"
    />
    <comment-component :commentArray="repliesArray" v-on:refresh-event="populateRepliesArray()"/>
  </div>
</template>

<script>
import Post from "@/components/Post.vue";
import API_Service from "@/service/API_Service.js";
import CommentComponent from "@/components/CommentComponent.vue";

export default {
  data() {
    return {
      parentPost: {
      },
      repliesArray: [],
      formVisible: false,
    };
  },
  components: {
    Post,
    CommentComponent
  },
  methods: {
    populateRepliesArray() {
      API_Service.getPostById(this.$route.params.subnoddit, this.$route.params.id)
      .then(returnedPost => (this.parentPost = returnedPost))
      .then(
        API_Service.getRepliesForPost(
          this.$route.params.subnoddit,
          this.$route.params.id)
      .then(parsedReplies => (this.repliesArray = parsedReplies))
      );
    }
  },
  created() {
    this.populateRepliesArray();
  },
  beforeCreate() {
    // if thread path contains non numeric characters kick back to subnoddit
    let id = this.$route.params.id;
    if (id.match(/^[0-9]+$/) === null) {
      this.$router.push('/n/' + this.$route.params.subnoddit);
    }
  },
  beforeUpdate() {
    // if parentPost.postId does not exist then kick back to subnoddit
    //if (this.parentPost.postId == 0 | this.parentPost == undefined | Object.entries(this.parentPost).length == 0) {
    //           this.$router.push('/n/' + this.$route.params.subnoddit);
    //        }
  }
};
</script>

<style>
</style>