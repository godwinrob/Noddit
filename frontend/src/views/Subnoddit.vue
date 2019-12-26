<template>
  <div class="subnoddit">
    <div class="subnoddit-header">
      <div class="sub-header">
        <h2>Welcome to /n/{{this.$route.params.subnoddit}}!</h2>
        <h5>{{this.subnoddit.subnodditDescription}}</h5>
      </div>
      <div class="btn-container">
        <div class="sort-container">
          <gb-button
            color="black"
            id="sort-recent"
            class="btn btn-primary"
            @click="filterToggle"
            v-show="filter"
          >Sort by: Recent</gb-button>
          <gb-button
            color="black"
            id="sort-popular"
            class="btn btn-primary"
            @click="filterToggle"
            v-show="!filter"
          >Sort by: Popular</gb-button>
        </div>
       
          <div class="fav-button">
            <gb-button
              color="black"
              id="addFavorite"
              class="btn btn-primary"
              @click="addToFavorites"
              v-show="getUsername() != null && !favoriteButtonToggle"
            >Add to Favorites</gb-button>
    
            <gb-button  
            color="red"
            id="removeFavorite"
            class="btn btn-danger"
            @click="removeFromFavorites"
            v-show="getUsername() != null && favoriteButtonToggle"
            >Remove from Favorites</gb-button>
          </div>
        </div>
        
    </div>
    
    <forum-posts :postsArray="filteredPostsArray" />
    <div class="shadow p-3 mb-5 bg-white rounded"></div>
  </div>
</template>

<script>
import auth from "@/auth.js";
import API_Service from "@/service/API_Service.js";
import ForumPosts from "@/components/ForumPosts.vue";

export default {
  name: "home",
  components: {
    ForumPosts
  },
  data() {
    return {
      user: auth.getUser(),
      subnoddit: [],
      usersFavoritesArray: [],
      popularPostsArray: [],
      recentPostsArray: [],
      filter: false,
      favorite: {
        username: this.getUsername(),
        subnodditName: this.$route.params.subnoddit
      }
    };
  },
  computed: {
    favoriteButtonToggle() {
      if (
        this.usersFavoritesArray.find(
          obj => obj.subnodditId == this.subnoddit.subnodditId
        )
      ) {
        return true;
      } else {
        return false;
      }
    },
    filteredPostsArray() {
      if (this.filter) {
        return this.popularPostsArray;
      }
      else {
        return this.recentPostsArray;
      }
      
    }
  },
  methods: {
    getFavoritesForUser() {
      API_Service.getFavoritesForUser(this.getUsername()).then(
        parsedFavorites => (this.usersFavoritesArray = parsedFavorites)
      );
    },
    addToFavorites() {
      API_Service.addFavoriteSubnoddit(this.favorite).then(
        this.getFavoritesForUser()
      ).then(window.location.reload()); //TODO: find a better way to do this rather than reloading the whole page
    },
    removeFromFavorites() {
      API_Service.removeFavoriteSubnoddit(this.subnoddit.subnodditId).then(
        this.getFavoritesForUser()
      ).then(window.location.reload()); //TODO: find a better way to do this rather than reloading the whole page
    },
    getUsername() {
      if (auth.getUser() == null) {
        return null;
      } else {
        return auth.getUser().sub;
      }
    },
    filterToggle() {
      this.filter = !this.filter;
    }
  },
  created() {
    if (this.getUsername() != null) {
      this.getFavoritesForUser();
    }
    API_Service.getPostsForSubnoddit(this.$route.params.subnoddit).then(
      parsedPosts => (this.recentPostsArray = parsedPosts)
    );
    API_Service.getSubnodditByName(this.$route.params.subnoddit).then(
      parsedSubnoddit => (this.subnoddit = parsedSubnoddit)
    );
    API_Service.getPostsForSubnodditPopular(this.$route.params.subnoddit).then(
      parsedPosts => (this.popularPostsArray = parsedPosts)
    );
  }
};
</script>

<style scoped>
.subnoddit-header {
  margin: 5px;
  margin-left: 10px;
   display: grid;
   max-width: 99%;
grid-template-columns: 72% 28%; 
/*justify-content:space-between ;*/

}
 .btn-container{
  display: flex;
  justify-content: flex-end;
}
.sort-container{
  display: grid;
  height: 2rem;

}
h2{
  color: rgb(112, 112, 112);
}
h5{
 color: rgb(112, 112, 112);
}

@media only screen and (max-width: 414px) {
  .btn-container{
    flex-direction: column;
    justify-content: flex-start;
  }
  .sort-container{
    height: unset;
  
  }
}



.fav-button {
 height: 3rem;
 width: auto;
 margin-left: 5px;
 
  /*display: flex;
  flex-direction: row;
  justify-content: flex-end;*/

}


</style>