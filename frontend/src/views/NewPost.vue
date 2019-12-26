<template>
  <div id="post-form">
    <div class="shadow-lg p-3 mb-5 bg-white rounded">
      <div id="new-post" class="text-center">
        <form class="form-create-post" @submit.prevent="createPost">
          <h1 class="h3 mb-3 font-weight-normal">Create New Post</h1>
          <select class="form-control" id="subnoddit-dropdown" v-model="post.subnodditId" required>
            <option value="0" hidden>Select a subnoddit...</option>
            <option
              v-for="subnoddit in subnodditArray"
              :key="subnoddit.subnodditId"
              :value="subnoddit.subnodditId"
            >{{subnoddit.subnodditName}}</option>
          </select>

          <input
            type="text"
            id="title"
            class="form-control"
            placeholder="Title"
            v-model="post.title"
            required
            autofocus
          />
          <input
            type="text"
            id="image"
            class="form-control"
            placeholder="Image Link"
            v-model="post.imageAddress"
            required
            autofocus
          />
          <textarea
            type="textarea"
            id="category"
            class="form-control"
            placeholder="Body"
            v-model="post.body"
            rows="3"
            required
            autofocus
          />
          <gb-button
            v-b-modal.modal-2
            block
            color="black"
            class="btn btn-lg btn-primary btn-block"
            type="submit"
            @click="createPost"
          >Submit Post</gb-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import API_Service from "@/service/API_Service.js";
import auth from "@/auth.js";

export default {
  data() {
    return {
      subnodditArray: [],
      post: {
        subnodditId: 0,
        username: auth.getUser().sub,
        title: "",
        imageAddress: "",
        body: ""
      },
      invalidCredentials: false
    };
  },
  computed: {
    subnodditName: function() {
      return this.subnodditArray.find(
        obj => obj.subnodditId == this.post.subnodditId
      ).subnodditName;
    }
  },
  created() {
    API_Service.getAllSubnoddits().then(
      parsedSubnoddits => (this.subnodditArray = parsedSubnoddits)
    );
  },
  methods: {
    createPost() {
      if (
        (this.post.imageAddress == "") |
        (this.post.imageAddress == undefined)
      ) {
        this.post.imageAddress = "https://i.imgur.com/MySsgmm.jpg";
      }
      if (this.post.subnodditId != 0) {
        API_Service.createNewPost(this.post).then(() =>
          this.$emit("refresh-event")
        );
        this.$router.push({
          name: "subnoddit",
          params: { subnoddit: this.subnodditName }
        });
      }
    }
  }
};
</script>

<style scoped>
#post-form {
  margin-top: 100px;
  display: flex;
  justify-content: center;
  max-width: 100%;
}
#new-post {
  padding: 10px;
  background-color: rgb(245, 245, 245);
  border: 1px solid rgb(197, 197, 197);
  width: 370px;
  max-width: 100%;
}
.form-control {
  margin: 3px;
}
</style>