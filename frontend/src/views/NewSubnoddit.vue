<template>
  <div id="post-form">
    <div class="shadow-lg p-3 mb-5 bg-white rounded">
      <div id="new-sn" class="text-center">
        <form class="form-create-post" @submit.prevent="createSubnoddit">
          <h1 class="h3 mb-3 font-weight-normal">Create New Subnoddit</h1>
          <input
            type="text"
            id="name"
            class="form-control"
            placeholder="Subnoddit Name"
            v-model="subnoddit.subnodditName"
            required
            autofocus
          />
          <textarea
            type="textarea"
            id="description"
            class="form-control"
            placeholder="Subnoddit Description"
            v-model="subnoddit.subnodditDescription"
            rows="3"
            required
            autofocus
          />
          <gb-button
            color="black"
            class="btn btn-lg btn-primary btn-block"
            type="submit"
            @click.prevent="createSubnoddit"
          >Create Subnoddit</gb-button>
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
      subnoddit: {
        username: auth.getUser().sub,
        subnodditName: "",
        subnodditDescription: ""
      },
      subnodditSplitName: ''
    };
  },
  methods: {
    createSubnoddit() {
      this.subnodditSplitName = this.subnoddit.subnodditName.replace(' ', '_')
      API_Service.createNewSubnoddit(this.subnoddit);
      this.$router.push("/n/" + this.subnodditSplitName);
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
#new-sn {
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