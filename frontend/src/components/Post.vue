<template>
  <div>
    <div v-show="!commentPage" class="post-block">
      <div class="arrows">
        <div>
          <icon v-bind:disabled="hasVoted" @click.prevent="upvote" class="upvote" name="arrow-up"></icon>
        </div>
        <div>{{postScore}}</div>
        <div>
          <icon
            v-bind:disabled="hasVoted"
            class="downvote"
            name="arrow-down"
            @click.prevent="downvote"
          ></icon>
        </div>
      </div>
      <div class="image">
        <router-link :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }">
          <img class="post-image" :src="imageAddress" />
        </router-link>
      </div>
      <div class="title">
        <router-link :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }">
          <h5 style="color: black">{{title}}</h5>
        </router-link>
        <p v-show="commentPage">{{body}}</p>
        <div class="post-info">
          submitted on {{createdDate}} by
          <a :href="'/u/' + username">u/{{username}}</a> to
          <a :href="'/n/' + subnodditName">n/{{subnodditName}}</a>
        </div>
        <div class="actions">
          <span v-show="!commentPage">
            <router-link
              :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }"
            >comments</router-link>
          </span>
          <span v-show="commentPage && getUsername() != null"><a
            href="#"
            @click="populatePostObject"
          >reply</a></span>
          <span v-show="canDelete"> | <a href="#" @click.prevent="deletePost">delete</a></span>
        </div>
      </div>
      <div class="reply-comment">
        <form @submit.prevent="addNewReply" v-show="formVisible">
          <div class="form-element">
            <label for="post">Message Text:</label>
            <gb-textarea v-model="post.body" id="post" required></gb-textarea>
          </div>
          <span class="subButton">
          <gb-button type="submit" :disabled="isFormValid === false" color="black">Submit</gb-button>
          </span>
          <gb-button @click="resetForm" color="red">Cancel</gb-button>
        </form>
      </div>
    </div>
    <div v-show="commentPage" class="post-block-comment">
      <div class="arrows">
        <div>
          <icon v-bind:disabled="hasVoted" @click.prevent="upvote" class="upvote" name="arrow-up"></icon>
        </div>
        <div>{{postScore}}</div>
        <div>
          <icon
            v-bind:disabled="hasVoted"
            class="downvote"
            name="arrow-down"
            @click.prevent="downvote"
          ></icon>
        </div>
      </div>
      <div class="image">
        <router-link :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }">
          <img class="post-image-comment" :src="imageAddress" />
        </router-link>
      </div>
      <div class="title">
        <router-link :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }">
          <h5 style="color: black">{{title}}</h5>
        </router-link>
        <p v-show="commentPage">{{body}}</p>
        <div class="post-info">
          submitted on {{createdDate}} by
          <a :href="'/u/' + username">u/{{username}}</a> to
          <a :href="'/n/' + subnodditName">n/{{subnodditName}}</a>
        </div>
        <div class="actions">
          <span v-show="!commentPage">
            <router-link
              :to="{ name: 'subnoddit-post', params: { subnoddit: this.subnodditName, id: this.postId } }"
            >comments</router-link>
          </span>
          <span v-show="commentPage && getUsername() != null"><a
            href="#"
            @click="populatePostObject"
          >reply</a></span>
          <span v-show="canDelete"> | <a href="#" @click.prevent="deletePost">delete</a></span>
        </div>
      </div>
      <div class="reply-comment">
        <form @submit.prevent="addNewReply" v-show="formVisible">
          <div class="form-element">
            <label for="post">Message Text:</label>
            <gb-textarea v-model="post.body" id="post" required></gb-textarea>
          </div>
          <span class="subButton">
          <gb-button type="submit" :disabled="isFormValid === false" color="black">Submit</gb-button>
          </span>
          <gb-button @click="resetForm" color="red">Cancel</gb-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import auth from "@/auth.js";
import Icon from "vue-awesome/components/Icon";
import API_Service from "@/service/API_Service.js";

export default {
  props: {
    postScore: Number,
    postId: Number,
    imageAddress: String,
    title: String,
    subnodditId: Number,
    subnodditName: String,
    username: String,
    createdDate: String,
    body: String,
    commentPage: Boolean
  },
  data() {
    return {
      user: auth.getUser(),
      formVisible: false,
      post: {
        body: ""
      },
      votes: [],
      hasVoted: true,
      currentVote: false,
      vote: {
        username: "",
        postId: 0,
        vote: ""
      },
      favorite: {
        username: "",
        postId: 0
      },
      moderators: []
    };
  },

  components: {
    Icon
  },
  computed: {
    isFormValid() {
      return Boolean(this.post.body);
    },
    canDelete() {
      if (this.user != null) {
        if ((auth.getUser().rol === "admin") |
           (auth.getUser().rol === "super_admin") |
           (auth.getUser().sub === this.username) |
           (this.checkModerator())
          ) {
          return true;
        }
      } 
      return false;
    }
  },
  methods: {
    populateVotePostId() {
      setTimeout(() => {
        if ((this.postId !== undefined) | (this.postId !== null)) {
          this.vote.postId = this.postId;
          this.favorite.postId = this.postId;
        } else {
          this.vote.postId = this.comment.postId;
          this.favorite.postId = this.comment.postId;
        }
        API_Service.getVotesForPost(this.vote.postId).then(
          parsedVotes => (this.votes = parsedVotes)
        );
      }, 500);
    },
    getUsername() {
      if (auth.getUser() == null) {
        return null;
      } else {
        return auth.getUser().sub;
      }
    },
    upvote() {
      if (this.vote.postId == undefined) {
        this.vote.postId = this.postId;
      }
      if (this.hasVoted === true) {
        this.checkVoteStatus();
        if (this.hasVoted === false & this.currentVote === false) {
          this.vote.vote = "upvote";
          this.postScore++;
          this.hasVoted = true;
          this.currentVote = true;
          API_Service.addVoteToPost(this.vote);
        }
      }
    },
    downvote() {
      if (this.vote.postId == undefined) {
        this.vote.postId = this.postId;
      }
      if (this.hasVoted === true) {
        this.checkVoteStatus();
        if (this.hasVoted === false & this.currentVote === false) {
          this.vote.vote = "downvote";
          this.postScore--;
          this.hasVoted = true;
          this.currentVote = true;
          API_Service.addVoteToPost(this.vote);
        }
      }
    },
    addFavorite() {
      console.log("favorite!" + this.postId);
      API_Service.addFavoritePost(this.favorite);
    },
    deletePost() {
      console.log("Delete!" + this.postId);
      if (
        (auth.getUser().rol === "admin") |
        (auth.getUser().rol === "super_admin") |
        (auth.getUser().sub === this.username) |
        (this.checkModerator())
      ) {
        API_Service.deletePost(this.postId).then( () => this.$emit('refresh-event'));
      } else {
        console.log("No!");
      }
    },
    populatePostObject() {
      this.post = {
        body: "",
        subnodditName: this.subnodditName,
        subnodditId: this.subnodditId,
        parentPostId: this.postId,
        username: this.getUsername(),
        topLevelId: this.postId
      };
      this.formVisible = true;
    },
    resetForm() {
      this.post = {
        body: "",
        subnodditName: this.subnodditName,
        subnodditId: this.subnodditId,
        parentPostId: this.postId,
        username: this.getUsername(),
        topLevelId: this.postId
      };
      this.formVisible = false;
    },
    addNewReply() {
      API_Service.createNewReply(this.post).then( () => this.$emit('refresh-event')).then(this.formVisible = false);
    },
    checkVoteStatus() {
      if (this.votes.length < 1 && this.user != null) {
        this.hasVoted = false;
      }
      for (let i = 0; i < this.votes.length; i++) {
        if (this.votes[i].username == this.getUsername()) {
          this.hasVoted = true;
          break;
        } else {
          this.hasVoted = false;
        }
      }
    },
    checkModerator() {
      for (let i =0; i < this.moderators.length; i++) {
        if (this.moderators[i].username == this.getUsername()) {
          return true;
        }
      }
      return false;
    }
  },
  created() {
    this.vote.username = this.getUsername();
    this.favorite.username = this.getUsername();
    this.populateVotePostId();
    API_Service.getModerators(this.subnodditName).then(
      parsedMods => (this.moderators = parsedMods)
    );
  }
};
</script>

<style scoped>
.post-block {
  display: grid;
  grid-template-columns: 0.25fr 0.5fr 5fr;
  grid-template-areas:
    " arrows image title "
    "  .      .    reply ";
  margin: 2px;
  border-radius: 5px;
  padding: 10px;
  background-color: rgb(245, 245, 245);
  border: 1px solid rgb(197, 197, 197);
}
.post-block-comment {
  display: grid;
  grid-template-columns: 0.25fr 1.5fr 3fr;
  grid-template-areas:
    " arrows image title "
    "  arrows image    reply ";
  margin: 2px;
  border-radius: 5px;
  padding: 10px;
  background-color: rgb(245, 245, 245);
  border: 1px solid rgb(197, 197, 197);
}
.arrows {
  grid-area: arrows;
  text-align: center;
}
.image {
  grid-area: image;
}
.subButton {
  margin-right: 0.5rem;
}
.reply-comment {
  grid-area: reply;
  max-width: 100%;
}
.title {
  grid-area: title;
}
.post-image {
  max-width: 100px;
  max-height: 80px;
  margin: 10px;
  border-radius: 2px;
}
.post-image-comment {
  max-width: 80%;
  max-height: auto;
  margin: 20px;
  border-radius: 2px;
}
.upvote {
  color: orange;
}
.downvote {
  color: blue;
}
.form-element{
  max-width: 75%;
}
</style>