<template>
  <div class="post-block">
    <div class="body">
      <div>{{comment.body}}</div>
      <div class="post-info">
        submitted on {{comment.createdDate}} by
        <a
          :href="'/u/' + comment.username"
        >u/{{comment.username}}</a>
      </div>
      <div class="actions">
        <span><a href="#" @click="formVisible = true">reply</a></span>
      </div>
    </div>
    <div class="reply-comment">
      <form @submit.prevent="addNewReply" v-show="formVisible">
        <div class="form-element">
          <label for="post">Message Text:</label>
          <gb-textarea v-model="post.body" id="post" required></gb-textarea>
        </div>
        <span class="subButton">
          <gb-button type="submit" class="subButton" :disabled="isFormValid === false" color="black">Submit</gb-button>
        </span>
        <span>
          <gb-button @click="formVisible = false" color="red">Cancel</gb-button>
        </span>
      </form>
    </div>
    <div class="arrows">
      <div>
        <icon class="upvote" name="arrow-up" @click.prevent="upvote"></icon>
      </div>
      <div>{{comment.postScore}}</div>
      <div>
        <icon class="downvote" name="arrow-down" @click.prevent="downvote"></icon>
      </div>
    </div>
  </div>
</template>

<script>
import Icon from "vue-awesome/components/Icon";
import auth from "../auth";
import API_Service from "@/service/API_Service.js";

export default {
  data() {
    return {
      user: auth.getUser(),
      formVisible: false,
      post: {
        body: "",
        subnodditName: this.comment.subnodditName,
        subnodditId: this.comment.subnodditId,
        parentPostId: this.comment.postId,
        username: this.getUsername,
        topLevelId: this.comment.topLevelId
      },
      votes: [],
      hasVoted: true,
      currentVote: false,
      vote: {
        username: this.getUsername(),
        postId: this.comment.postId,
        vote: ""
      }
    };
  },
  computed: {
    isFormValid() {
      return Boolean(this.post.body);
    }
  },
  props: {
    comment: Object,
    postId: Number
  },
  components: {
    Icon
  },
  methods: {
    resetForm() {
      this.post = {
        body: "",
        subnodditName: this.comment.subnodditName,
        subnodditId: this.comment.subnodditId,
        parentPostId: this.comment.postId,
        username: this.getUsername,
        topLevelId: this.comment.topLevelId
      };
      this.formVisible = false;
    },
    addNewReply() {
      API_Service.createNewReply(
        this.post,
        this.$route.params.subnoddit,
        this.$route.params.id
      ).then( () => this.$emit('refresh-event')).then(this.formVisible = false);
    },
    getUsername() {
      if (auth.getUser() == null) {
        return null;
      } else {
        return auth.getUser().sub;
      }
    },
    upvote() {
      if (this.hasVoted === true) {
        this.checkVoteStatus();
        if (this.hasVoted === false & this.currentVote === false) {
          this.vote.vote = "upvote";
          this.comment.postScore++;
          this.hasVoted = true;
          this.currentVote = true;
          API_Service.addVoteToPost(this.vote);
        }
      }
    },
    downvote() {
      if (this.hasVoted === true) {
        this.checkVoteStatus();
        if (this.hasVoted === false & this.currentVote === false) {
          this.vote.vote = "downvote";
          this.comment.postScore--;
          this.hasVoted = true;
          this.currentVote = true;
          API_Service.addVoteToPost(this.vote);
        }
      }
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
    }
  },
  created() {
    API_Service.getVotesForPost(this.comment.postId).then(
      parsedVotes => (this.votes = parsedVotes)
    );
  }
};
</script>

<style>

.post-block {
  box-sizing: border-box;
  display: grid;
  grid-template-columns: 0.25fr 5fr;
  grid-template-areas:
    " arrows body "
    "   .    reply ";
  margin: 2px;
  border-radius: 5px;
  padding: 10px;
  background-color: rgb(245, 245, 245);
  border: 1px solid rgb(197, 197, 197);
}
.arrows {
  box-sizing: border-box;
  grid-area: arrows;
  text-align: center;
  margin-right: 4px;
}
.body {
  box-sizing: border-box;
  grid-area: body;
  resize: vertical;
}
.reply-comment {
  box-sizing: border-box;
  grid-area: reply;
  
}
.post-info {
  box-sizing: border-box;
  margin-top: 15px;
  font-size: 0.7em;
  resize: vertical;
}
.actions {
  box-sizing: border-box;
  font-size: 0.7em;
}
.post-image {
  box-sizing: border-box;
  max-width: 100px;
  max-height: 80px;
  margin: 10px;
  border-radius: 2px;
  resize: vertical;
}
.upvote {
  box-sizing: border-box;
  color: orange;
}
.downvote {
  box-sizing: border-box;
  color: blue;
}

.subButton {
  box-sizing: border-box;
  margin-right: 0.25rem;
}

#post {
  box-sizing: border-box;
  max-width: 100%;
  height: 75%;
  outline: 0;
  
}
.title{
  resize: vertical;
}
.form-element{
  max-width: 50%
}


</style>