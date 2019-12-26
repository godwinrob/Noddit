import auth from '../auth.js'

export default {
  getRecentPosts() {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/recentposts`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  getAllSubnoddits() {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/subnoddits`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  getPostById(subnoddit, postId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/${subnoddit}/${postId}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  createNewPost(post) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/post/create`, {
      method: "POST",
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify(
        post
      )
    }).catch(err => console.error(err));
  },
  getPostsForSubnoddit(subnoddit) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/allposts/${subnoddit}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  getPostsForSubnodditPopular(subnoddit) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/allpostspopular/${subnoddit}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  getRepliesForPost(subnoddit, postId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/${subnoddit}/${postId}/replies`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  getFavoritesForUser(user) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/favorites/${user}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
      .then((response) => {
        return response.json();
      });
  },
  deletePost(postId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/post/delete/${postId}`, {
      method: 'DELETE',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + auth.getToken()
      }
    });
  },
  createNewSubnoddit(subnoddit) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/subnoddits/create`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + auth.getToken()
      },
      body: JSON.stringify(
        subnoddit
      )
    }).catch(err => console.error(err));
  },
  createNewReply(post) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/${post.subnodditName}/${post.topLevelId}/createreply`, {
      method: "POST",
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        body: post.body,
        subnodditName: post.subnodditName,
        subnodditId: post.subnodditId,
        parentPostId: post.parentPostId,
        username: post.username,
        topLevelId: post.topLevelId
      })
    }).catch(err => console.error(err));
  },
  getSubnodditByName(subnodditName) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/subnoddits/${subnodditName}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  getVotesForPost(postId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/post/votes/${postId}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  addVoteToPost(vote) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/post/vote`, {
      method: 'POST',
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify(
        vote
      )
    }).catch(err => console.error(err));
  },
  addFavoriteSubnoddit(favorite) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/favorites/create/subnoddit`, {
      method: 'POST',
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify(
        favorite
      )
    }).catch(err => console.error(err));
  },
  removeFavoriteSubnoddit(subnodditId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/favorites/subnoddit/${subnodditId}`, {
      method: 'DELETE',
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      }
    }).catch(err => console.error(err));
  },
  addFavoritePost(favorite) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/favorites/create/post`, {
      method: 'POST',
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify(
        favorite
      )
    }).catch(err => console.error(err));
  },
  getActiveSubnoddits() {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/subnoddits/active`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  getPopularPosts() {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/popularposts`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  searchSubnoddits(searchTerm) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/subnoddits/search/${searchTerm}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  getModerators(subnodditId) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/moderators/${subnodditId}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  getPostsForUser(username) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/public/post/user/${username}`, { headers: { Authorization: 'Bearer ' + auth.getToken() } })
    .then((response) => {
      return response.json();
    });
  },
  updateUsername(user) {
    return fetch(`${process.env.VUE_APP_REMOTE_API}/api/user/update/username/${user.username}`, {
      method: "PUT",
      headers: {
        Authorization: "Bearer " + auth.getToken(),
        "Content-Type": "application/json"
      },
      body: JSON.stringify(
        user
      )
    }).catch(err => console.error(err));
  }
  
}