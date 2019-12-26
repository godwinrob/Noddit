<template>
  <div class="search-subnoddits">
    <div class="header-div">
      <h1>Search Results</h1>
    </div>
    <b-list-group>
      <b-list-group-item v-for="subnoddit in searchArray" :key="subnoddit.subnodditId">
        <b-link :to="{ name: 'subnoddit', params: { subnoddit: subnoddit.subnodditName } }">
          <div>/n/{{subnoddit.subnodditName}} - {{subnoddit.subnodditDescription}}</div>
        </b-link>
      </b-list-group-item>
      <h4 class="noresults" v-show="emptyArray">No Results Found for Search "{{this.searchTerm}}"</h4>
    </b-list-group>
  </div>
</template>

<script>
import auth from "@/auth.js";
import API_Service from "@/service/API_Service.js";

export default {
  data() {
    return {
      user: auth.getUser(),
      searchArray: [],
      searchTerm: this.$route.params.searchTerm
    };
  },
  created() {
    if (this.searchTerm != undefined) {
      API_Service.searchSubnoddits(this.searchTerm).then(
        searchResults => (this.searchArray = searchResults)
      );
    }
  },
  computed: {
    emptyArray() {
      if (this.searchArray.length < 1 || this.searchTerm == undefined) {
        return true;
      } else {
        return false;
      }
    }
  }
};
</script>

<style scoped>
ul {
  list-style-type: none;
}
.header-div {
  margin: 5px;
}
.li-div {
  max-width: 80%;
  width: 600px;
  padding: 3px;
  border: 3px solid rgb(190, 189, 189);
  margin: 5px;
  border-radius: 4px;
  background-color: rgb(216, 211, 211);
}
.no-results {
  max-width: 90%;
  margin: 5px;
}
</style>